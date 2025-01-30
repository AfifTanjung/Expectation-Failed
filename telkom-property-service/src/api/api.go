package api

import (
	"log"
	"telkom-property-service/pkg"
	alatController "telkom-property-service/src/feature/alat/controller"
	alatRepo "telkom-property-service/src/feature/alat/repository"
	alatService "telkom-property-service/src/feature/alat/services"

	RawController "telkom-property-service/src/feature/raw/controller"
	rawRepo "telkom-property-service/src/feature/raw/repository"
	rawService "telkom-property-service/src/feature/raw/service"

	validCont "telkom-property-service/src/feature/valid/controller"
	validRepo "telkom-property-service/src/feature/valid/repository"
	validServ "telkom-property-service/src/feature/valid/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func SetupRoutes() *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}
	db, err := pkg.InitDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	logger := log.Default()

	alatRepo := alatRepo.NewAlatRepository(db, logger)
	alatService := alatService.NewAlatService(alatRepo, logger)
	alatCtrl := alatController.NewAlatController(alatService)

	rawRepo := rawRepo.NewRawRepository(db, logger)
	rawService := rawService.NewRawService(rawRepo, logger)
	rawCtrl := RawController.NewRawController(rawService)

	validRepo := validRepo.NewValidationRepository(db, logger)
	validService := validServ.NewValidationService(validRepo)
	validCtrl := validCont.NewValidationController(validService)

	r := gin.Default()

	r.GET("/home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "HTTP server listening",
			"success": true,
		})
	})

	r.GET("/server-time", func(ctx *gin.Context) {
		location, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": "Failed to load timezone",
			})
			return
		}
		currentTime := time.Now().In(location)
		serverTime := currentTime.Format("2006-01-02 15:04:05")
		ctx.JSON(200, gin.H{
			"server_time": serverTime,
			"timezone":    location.String(),
		})
	})

	alatsGroup := r.Group("/alats")
	{
		alatsGroup.POST("", alatCtrl.CreateAlat)
		alatsGroup.GET("/lantai/:lantai", alatCtrl.GetByLantai)
		alatsGroup.GET("/stasiun/:stasiun", alatCtrl.GetByStasiun)
	}

	rawsGroup := r.Group("/raws")
	{
		rawsGroup.POST("/:alatId", rawCtrl.CreateByAlatId)
		rawsGroup.GET("/:alatId", rawCtrl.GetById)
	}

	validationGroup := r.Group("/validations")
	{
		validationGroup.POST("/validate", validCtrl.ValidateDataHandler)
	}

	return r
}
