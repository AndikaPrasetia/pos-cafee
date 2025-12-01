package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AndikaPrasetia/pos-cafee/internal/cache"
	"github.com/AndikaPrasetia/pos-cafee/internal/config"
	"github.com/AndikaPrasetia/pos-cafee/internal/handlers"
	"github.com/AndikaPrasetia/pos-cafee/internal/middleware"
	"github.com/AndikaPrasetia/pos-cafee/internal/repositories"
	"github.com/AndikaPrasetia/pos-cafee/internal/services"
	"github.com/AndikaPrasetia/pos-cafee/pkg/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connection
	db := config.ConnectDB(cfg)
	defer config.CloseDB(db)

	// Initialize Redis connection
	rdb := config.RedisClient(cfg)
	defer config.CloseRedis(rdb)

	// Initialize cache
	cacheClient := cache.NewRedisCache(rdb)

	// Initialize repositories
	repo := repositories.NewRepository(db)

	// Initialize services
	authService := services.NewAuthService(repo.UserRepo, cfg.JWTSecret, parseDuration(cfg.JWTExpiry))
	menuService := services.NewMenuService(repo.MenuRepo, repo.InventoryRepo, cacheClient)
	orderService := services.NewOrderService(repo.OrderRepo, repo.OrderItemRepo, repo.MenuRepo, repo.InventoryRepo, repo.StockTransactionRepo, cacheClient)
	inventoryService := services.NewInventoryService(repo.InventoryRepo, repo.StockTransactionRepo, repo.MenuRepo)
	expenseService := services.NewExpenseService(repo.ExpenseRepo)
	reportService := services.NewReportService(repo.OrderRepo, repo.MenuRepo, repo.InventoryRepo, repo.ExpenseRepo, repo.Queries, cacheClient)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	menuHandler := handlers.NewMenuHandler(menuService)
	orderHandler := handlers.NewOrderHandler(orderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	reportHandler := handlers.NewReportHandler(reportService)

	// Initialize Gin router
	router := gin.New()

	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Add request size limit (e.g., 8MB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// Add basic health check endpoint - we'll use the maintenance handler
	maintenanceHandler := handlers.NewMaintenanceHandler()
	router.GET("/health", maintenanceHandler.HealthCheck)

	// Public routes (no authentication required)
	public := router.Group("/api/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
	}

	// Authentication protected routes (authentication required)
	authProtected := router.Group("/api/auth")
	authProtected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		authProtected.GET("/profile", authHandler.Profile)
		authProtected.PUT("/change-password", authHandler.ChangePassword)
		authProtected.POST("/logout", authHandler.Logout)
	}

	// Menu management routes (require manager or admin role)
	menu := router.Group("/api/menu")
	menu.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "manager"))
	{
		// Category endpoints
		menu.GET("/categories", menuHandler.ListCategories)
		menu.POST("/categories", menuHandler.CreateCategory)
		menu.GET("/categories/:id", menuHandler.GetCategory)
		menu.PUT("/categories/:id", menuHandler.UpdateCategory)
		menu.DELETE("/categories/:id", menuHandler.DeleteCategory)

		// Menu item endpoints
		menu.GET("/items", menuHandler.ListMenuItems)
		menu.POST("/items", menuHandler.CreateMenuItem)
		menu.GET("/items/:id", menuHandler.GetMenuItem)
		menu.PUT("/items/:id", menuHandler.UpdateMenuItem)
		menu.DELETE("/items/:id", menuHandler.DeleteMenuItem)
	}

	// Order management routes (require cashier role or higher)
	orders := router.Group("/api/orders")
	orders.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "cashier"))
	{
		orders.GET("/", orderHandler.ListOrders)
		orders.POST("/", orderHandler.CreateOrder)
		orders.GET("/:id", orderHandler.GetOrder)
		orders.POST("/:id/items", orderHandler.AddItemToOrder)
		orders.PUT("/:id/complete", orderHandler.CompleteOrder)
		orders.PUT("/:id/cancel", orderHandler.CancelOrder)
	}

	// Inventory management routes (require manager or admin role)
	inventory := router.Group("/api/inventory")
	inventory.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "manager"))
	{
		inventory.GET("/", inventoryHandler.ListInventory)
		inventory.GET("/low-stock", inventoryHandler.GetLowStockItems)
		inventory.POST("/adjust", inventoryHandler.UpdateInventory)
		inventory.GET("/transactions", inventoryHandler.ListStockTransactions)
	}

	// Reporting routes (require manager or admin role)
	reports := router.Group("/api/reports")
	reports.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "manager"))
	{
		reports.GET("/daily-sales", reportHandler.GetDailySalesReport)
		reports.GET("/financial-summary", reportHandler.GetFinancialSummaryReport)
		reports.GET("/sales-by-category", reportHandler.GetSalesByCategoryReport)
		reports.GET("/top-selling-items", reportHandler.GetTopSellingItemsReport)
	}

	// Expense management routes (require manager or admin role)
	expenses := router.Group("/api/expenses")
	expenses.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "manager"))
	{
		expenses.GET("/", expenseHandler.ListExpenses)
		expenses.GET("/summary", expenseHandler.GetExpenseSummary)
		expenses.POST("/", expenseHandler.CreateExpense)
		expenses.GET("/:id", expenseHandler.GetExpense)
		expenses.PUT("/:id", expenseHandler.UpdateExpense)
		expenses.DELETE("/:id", expenseHandler.DeleteExpense)
	}

	// Maintenance routes (admin for backups)
	// Note: Health check already added earlier
	maintenance := router.Group("/api/maintenance")
	maintenance.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "admin"))
	{
		maintenance.POST("/backup", maintenanceHandler.DatabaseBackup)
	}

	// Create HTTP server with timeout settings
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start the server
	log.Printf("Starting server on port %s in %s mode", cfg.Port, cfg.Environment)
	log.Fatal(srv.ListenAndServe())
}

// parseDuration parses the duration string from config
func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Default to 24 hours if parsing fails
		return 24 * time.Hour
	}
	return duration
}
