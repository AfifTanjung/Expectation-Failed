package controller

import (
	"context"
	"errors"
	"net/http"
	"telkom-property-service/src/feature/valid/service"
	"time"

	"github.com/gin-gonic/gin"
)

type ValidationController struct {
	service service.ValidationServices
}

func NewValidationController(service service.ValidationServices) *ValidationController {
	return &ValidationController{service: service}
}

func (c *ValidationController) ValidateDataHandler(ctx *gin.Context) {
	requestCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Minute)
	defer cancel()
	err := c.service.ValidateData(requestCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			ctx.JSON(http.StatusRequestTimeout, gin.H{
				"status":  "error",
				"message": "Request timed out or was cancelled.",
				"details": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Internal server error occurred during validation.",
			"details": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Validation completed successfully.",
	})
}
