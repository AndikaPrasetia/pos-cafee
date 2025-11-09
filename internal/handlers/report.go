package handlers

import (
	"net/http"
	"strconv"

	"github.com/AndikaPrasetia/pos-cafee/internal/services"
	"github.com/AndikaPrasetia/pos-cafee/pkg/types"
	"github.com/gin-gonic/gin"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	reportService *services.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// GetDailySalesReport handles daily sales report requests
func (h *ReportHandler) GetDailySalesReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("date parameter is required in YYYY-MM-DD format"))
		return
	}

	result, err := h.reportService.GetDailySalesReport(dateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetFinancialSummaryReport handles financial summary report requests
func (h *ReportHandler) GetFinancialSummaryReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("start_date and end_date parameters are required in YYYY-MM-DD format"))
		return
	}

	result, err := h.reportService.GetFinancialSummaryReport(startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSalesByCategoryReport handles sales by category report requests
func (h *ReportHandler) GetSalesByCategoryReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("start_date and end_date parameters are required in YYYY-MM-DD format"))
		return
	}

	result, err := h.reportService.GetSalesByCategoryReport(startDateStr, endDateStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetTopSellingItemsReport handles top selling items report requests
func (h *ReportHandler) GetTopSellingItemsReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, types.APIResponseWithError("start_date and end_date parameters are required in YYYY-MM-DD format"))
		return
	}

	// Get limit parameter (optional)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	result, err := h.reportService.GetTopSellingItemsReport(startDateStr, endDateStr, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponseWithError(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}