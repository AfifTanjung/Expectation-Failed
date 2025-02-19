package controller

import (
	"net/http"
	"strconv"
	"telkom-property-service/src/feature/maintenance/service"
	"telkom-property-service/src/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MaintenanceController struct {
	service service.MaintenanceService
}

func NewMaintenanceController(service service.MaintenanceService) *MaintenanceController {
	return &MaintenanceController{service: service}
}

func (mc *MaintenanceController) Create(c *gin.Context) {
	var maintenance model.Maintenance
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := mc.service.Create(c.Request.Context(), &maintenance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create maintenance log", "detail":err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":          true,
		"message":         "Maintenance log created successfully",
		"Log Maintenance": maintenance,
	})
}

func (mc *MaintenanceController) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var maintenance model.Maintenance
	if err := c.ShouldBindJSON(&maintenance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := mc.service.Update(c.Request.Context(), id, &maintenance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update maintenance log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"message":    "Maintenance log updated successfully",
		"update Log": maintenance,
	})
}

func (mc *MaintenanceController) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := mc.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete maintenance log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Maintenance log deleted successfully",
	})
}

func (mc *MaintenanceController) GetLogById(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format", "detail": err.Error()})
		return
	}

	maintenance, err := mc.service.GetLogById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance log not found"})
		return
	}

	 if maintenance.ID == uuid.Nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Maintenance log not found"})
        return
    }


	c.JSON(http.StatusOK, gin.H{
		"status":  false,
		"message": "succesfully Fetching data",
		"data":    maintenance,
	})
}

func (mc *MaintenanceController) GetLogByLantai(c *gin.Context) {
	lantai, err := strconv.Atoi(c.Param("lantai"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lantai format"})
		return
	}

	logs, err := mc.service.GetLogByLantai(c.Request.Context(), lantai)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "succesfully fetch data ",
		"data":    logs,
	})
}
