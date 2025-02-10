package controller

import (
	"net/http"
	"telkom-property-service/src/feature/raw/service"
	"telkom-property-service/src/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RawController struct {
	service service.RawServices
}

func NewRawController(service service.RawServices) *RawController {
	return &RawController{service: service}
}

func (c *RawController) CreateByAlatId(ctx *gin.Context) {
	var raw model.Raw
	if err := ctx.ShouldBindJSON(&raw); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid input: " + err.Error(),
		})
		return
	}
	if raw.AlatID == uuid.Nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "AlatID tidak boleh kosong",
		})
		return
	}
	if err := c.service.CreateByAlatId(ctx.Request.Context(), raw); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to create Raw data: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"dataRaw": gin.H{
			"id":          raw.ID,
			"waktuMasuk":  raw.WaktuMasuk,
			"WaktuTerima": raw.WaktuTerima,
			"co2":         raw.Co2,
			"pm2_5":       raw.Pm2_5,
			"temperature": raw.Temp,
			"humidity":    raw.Hum,
			"windSpeed":   raw.Ws,
			"windDirec":   raw.Wd,
			"alatId":      raw.AlatID.String(),
		},
	})
}

func formatRawData(raws []model.Raw) []gin.H {
	var formattedData []gin.H
	for _, raw := range raws {
		formattedData = append(formattedData, gin.H{
			"id":          raw.ID,
			"waktuMasuk":  raw.WaktuMasuk,
			"WaktuTerima": raw.WaktuTerima,
			"co2":         raw.Co2,
			"pm2_5":       raw.Pm2_5,
			"temperature": raw.Temp,
			"humidity":    raw.Hum,
			"windSpeed":   raw.Ws,
			"windDirec":   raw.Wd,
		})
	}
	return formattedData
}

func (c *RawController) GetById(ctx *gin.Context) {
	alatID := ctx.Param("alatId")
	parsedAlatID, err := uuid.Parse(alatID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Invalid alat ID format",
		})
		return
	}
	alat, rawData, err := c.service.GetById(ctx.Request.Context(), parsedAlatID)
	if err != nil {
		if err == service.ErrAlatNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"status": "error",
				"error":  "Alat not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
			"error":  "Failed to fetch alat or raw data: " + err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"alat_uuid": alat.ID.String(),
			"namaAlat":  alat.Nama,
			"type":      alat.Type,
			"dataRaw":   formatRawData(rawData),
		},
	})
}
