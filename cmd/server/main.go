package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	// Load configuration first to get environment info
	cfg := config.LoadConfig()

	// Initialize logger with environment and log level from config
	utils.InitLogger(cfg.Environment, cfg.LogLevel)

	// Log application startup
	utils.LogInfo("Starting POS Cafe server", map[string]any{
		"environment": cfg.Environment,
		"port":        cfg.Port,
		"version":     "3.0.0",
	})

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Initialize database connection
	db := config.ConnectDB(cfg)

	// Initialize Redis manager with automatic reconnection
	redisManager, err := cache.NewRedisManager(cfg)
	var cacheClient cache.Cache

	if err != nil {
		log.Printf("Warning: Failed to initialize Redis connection: %v", err)
		log.Println("Using fallback cache only. The application will continue to run without Redis.")

		// Initialize adaptive cache with a nil Redis manager (will use fallback only)
		// We need to modify the NewAdaptiveCache function to handle nil RedisManager
		cacheClient = cache.NewAdaptiveCache(nil)  // Will use fallback cache only initially
	} else {
		// Initialize adaptive cache with Redis and fallback
		cacheClient = cache.NewAdaptiveCache(redisManager)

		// Start Redis monitoring in the background
		go redisManager.Monitor(context.Background())
	}

	// Initialize repositories
	repo := repositories.NewRepository(db)

	// Initialize services
	authService := services.NewAuthService(repo.UserRepo, cfg.JWTSecret, parseDuration(cfg.JWTExpiry))
	menuService := services.NewMenuService(repo.MenuRepo, repo.InventoryRepo, cacheClient)
	orderService := services.NewOrderService(repo.OrderRepo, repo.OrderItemRepo, repo.MenuRepo, repo.InventoryRepo, repo.StockTransactionRepo, cacheClient)
	inventoryService := services.NewInventoryService(repo.InventoryRepo, repo.StockTransactionRepo, repo.MenuRepo)
	expenseService := services.NewExpenseService(repo.ExpenseRepo)
	reportService := services.NewReportService(repo.OrderRepo, repo.MenuRepo, repo.InventoryRepo, repo.ExpenseRepo, repo.Queries, cacheClient)
	userService := services.NewUserService(repo.UserRepo)

	// Initialize cache warming service and warm cache in the background
	cacheWarmingService := cache.NewCacheWarmingService(cacheClient, repo)
	cacheWarmingService.WarmCacheInBackground()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	menuHandler := handlers.NewMenuHandler(menuService)
	orderHandler := handlers.NewOrderHandler(orderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	reportHandler := handlers.NewReportHandler(reportService)
	userHandler := handlers.NewUserHandler(authService, userService)  // Note: userService is needed here
	systemHandler := handlers.NewSystemHandler(redisManager)

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

	// User management routes (admin only)
	users := router.Group("/api/users")
	users.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "admin"))
	{
		users.GET("/", userHandler.ListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.PUT("/:id/deactivate", userHandler.DeactivateUser)
		users.PUT("/:id/activate", userHandler.ActivateUser)
	}

	// System monitoring routes (admin only)
	system := router.Group("/api/system")
	system.Use(middleware.RoleAuthMiddleware(cfg.JWTSecret, "admin"))
	{
		system.GET("/redis/health", systemHandler.GetRedisHealth)
		system.GET("/redis/metrics", systemHandler.GetRedisMetrics)
	}

	// Create HTTP server with timeout settings
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// channel for signal interrupt
	quit := make(chan os.Signal, 1)
	// signal interrupted (Ctrl+C)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// run server in goroutine
	go func() {
	log.Printf("Starting server on port %s in %s mode", cfg.Port, cfg.Environment)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start server: %v", err))
		}
	}()

	// blocking main goroutine until signal recieved
	<-quit
	fmt.Println("\nShutting down server...")

	// timeout 5 seconds for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
	}

	// close db connection
	if err := db.Close(); err != nil {
		log.Printf("Error closing database: %v\n", err)
	}

	// close adaptive cache if it supports closing
	if closer, ok := cacheClient.(interface{ Close() }); ok {
		closer.Close()
	}

	// close redis connection if Redis manager is available
	if redisManager != nil {
		if err := redisManager.GetClient().Close(); err != nil {
			log.Printf("Error closing redis: %v\n", err)
		}
	}

	fmt.Println("Server gracefully stopped 󱠡 ")
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
