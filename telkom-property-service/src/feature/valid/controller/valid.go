package controller

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"telkom-property-service/src/feature/valid/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ValidationController struct {
	service service.ValidationServices
	logger  *logrus.Logger
}

type ValidationResponse struct {
	Status    string      `json:"status"`
	Message   string      `json:"message"`
	Details   string      `json:"details,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewValidationController(service service.ValidationServices, logger *logrus.Logger) *ValidationController {
	return &ValidationController{
		service: service,
		logger:  logger,
	}
}

func (c *ValidationController) ValidateDataHandler(ctx *gin.Context) {
	c.handleValidation(ctx, false, 0)
}

func (c *ValidationController) ValidateBatchHandler(ctx *gin.Context) {
	c.handleValidation(ctx, true, getBatchSize(ctx))
}

func (c *ValidationController) handleValidation(ctx *gin.Context, isBatch bool, batchSize int) {
	requestID := ctx.GetString("request_id")
	c.logger.WithFields(logrus.Fields{
		"request_id": requestID,
		"is_batch":   isBatch,
		"batch_size": batchSize,
	}).Info("Starting validation request")

	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Minute)
	defer cancel()

	// Start time for request duration tracking
	startTime := time.Now()

	var err error
	if isBatch {
		err = c.service.ValidateBatch(requestCtx, batchSize)
	} else {
		err = c.service.ValidateData(requestCtx)
	}

	duration := time.Since(startTime)

	if err != nil {
		c.handleError(ctx, err, requestID, duration)
		return
	}

	response := ValidationResponse{
		Status:    "success",
		Message:   "Validation completed successfully",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"duration_ms": duration.Milliseconds(),
			"is_batch":    isBatch,
			"batch_size":  batchSize,
		},
	}

	c.logger.WithFields(logrus.Fields{
		"request_id":  requestID,
		"duration_ms": duration.Milliseconds(),
		"status":      "success",
	}).Info("Validation completed")

	ctx.JSON(http.StatusOK, response)
}

func (c *ValidationController) handleError(ctx *gin.Context, err error, requestID string, duration time.Duration) {
	var response ValidationResponse
	var statusCode int

	switch {
	case errors.Is(err, context.DeadlineExceeded):
		statusCode = http.StatusRequestTimeout
		response = ValidationResponse{
			Status:  "error",
			Message: "Request timed out",
			Details: "The validation process took too long to complete",
		}
	case errors.Is(err, context.Canceled):
		statusCode = http.StatusBadRequest
		response = ValidationResponse{
			Status:  "error",
			Message: "Request was cancelled",
			Details: "The validation process was cancelled by the client",
		}
	case errors.Is(err, service.ErrValidation):
		statusCode = http.StatusUnprocessableEntity
		response = ValidationResponse{
			Status:  "error",
			Message: "Validation failed",
			Details: err.Error(),
		}
	default:
		statusCode = http.StatusInternalServerError
		response = ValidationResponse{
			Status:  "error",
			Message: "Internal server error occurred during validation",
			Details: "An unexpected error occurred while processing the request",
		}
	}

	response.Timestamp = time.Now()
	response.Data = map[string]interface{}{
		"duration_ms": duration.Milliseconds(),
	}

	c.logger.WithFields(logrus.Fields{
		"request_id":  requestID,
		"error":       err.Error(),
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}).Error("Validation failed")

	ctx.JSON(statusCode, response)
}

func getBatchSize(ctx *gin.Context) int {
	batchSizeStr := ctx.DefaultQuery("batch_size", "100")
	batchSize, err := strconv.Atoi(batchSizeStr)
	if err != nil || batchSize <= 0 {
		return 100
	}
	return batchSize
}

func (c *ValidationController) HealthCheckHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ValidationResponse{
		Status:    "success",
		Message:   "Validation service is healthy",
		Timestamp: time.Now(),
	})
}
