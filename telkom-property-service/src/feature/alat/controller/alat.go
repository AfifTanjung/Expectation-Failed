package controller

import (
	"net/http"
	"strconv"
	"telkom-property-service/src/feature/alat/services"
	"telkom-property-service/src/model"

	"github.com/gin-gonic/gin"
)

type AlatController struct {
	service services.AlatService
}

func NewAlatController(service services.AlatService) *AlatController {
	return &AlatController{service: service}
}

func (c *AlatController) CreateAlat(ctx *gin.Context) {
	var alat model.Alat
	if err := ctx.ShouldBindJSON(&alat); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid input",
			"error":   err.Error(),
		})
		return
	}

	if err := c.service.CreateAlat(ctx.Request.Context(), alat); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create alat",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Alat created successfully",
		"data":    alat,
	})
}

func (c *AlatController) GetByLantai(ctx *gin.Context) {
	lantaiParam := ctx.Param("lantai")
	lantai, err := strconv.Atoi(lantaiParam)
	if err != nil || lantai <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid lantai parameter, must be a positive number",
			"error":   err.Error(),
		})
		return
	}
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if limit < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Limit must be a non-negative number",
		})
		return
	}
	if offset < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Offset must be a non-negative number",
		})
		return
	}
	alats, err := c.service.GetAlatByLantai(ctx.Request.Context(), lantai, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch alat",
			"error":   err.Error(),
		})
		return
	}

	if len(alats) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "No alat found for the specified lantai",
			"data":    []model.Alat{},
		})
		return
	}

	responseData := []gin.H{}
	for _, alat := range alats {
		responseData = append(responseData, gin.H{
			"id":         alat.ID,
			"nama":       alat.Nama,
			"type":       alat.Type,
			"lantai":     alat.Lantai,
			"stasiun":    alat.Stasiun,
			"created_at": alat.CreatedAt,
			"updated_at": alat.UpdatedAt,
			"raw_count":  len(alat.Raws),
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alat fetched successfully by lantai",
		"data":    responseData,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
			"total":  len(alats),
		},
	})
}

func (c *AlatController) GetByStasiun(ctx *gin.Context) {
	stasiun := ctx.Param("stasiun")
	if stasiun == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Stasiun parameter is required",
		})
		return
	}
	alats, err := c.service.GetAlatByStasiun(ctx.Request.Context(), stasiun)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch alat",
			"error":   err.Error(),
		})
		return
	}
	if len(alats) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": true,
			"message": "No alat found for the specified stasiun",
			"data":    []model.Alat{},
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Alat fetched successfully by stasiun",
		"data":    alats,
	})
}
