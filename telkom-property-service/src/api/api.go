package api

import (
	"log"
	"time"

	"telkom-property-service/src/feature/admin/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	alatController "telkom-property-service/src/feature/alat/controller"
	alatRepo "telkom-property-service/src/feature/alat/repository"
	alatService "telkom-property-service/src/feature/alat/services"

	rawController "telkom-property-service/src/feature/raw/controller"
	rawRepo "telkom-property-service/src/feature/raw/repository"
	rawService "telkom-property-service/src/feature/raw/service"

	validCont "telkom-property-service/src/feature/valid/controller"
	validRepo "telkom-property-service/src/feature/valid/repository"
	validServ "telkom-property-service/src/feature/valid/service"

	maintController "telkom-property-service/src/feature/maintenance/controller"
	maintRepo "telkom-property-service/src/feature/maintenance/repository"
	maintService "telkom-property-service/src/feature/maintenance/service"

	authController "telkom-property-service/src/feature/admin/controller"
	authRepo "telkom-property-service/src/feature/admin/repository"
	authService "telkom-property-service/src/feature/admin/service"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	// Pastikan db tidak nil sebelum digunakan
	if db == nil {
		log.Fatal("Database connection is nil, cannot proceed with route setup")
	}

	logger := logrus.New()

	// Auth Dependencies
	authRepository := authRepo.NewAuthRepository(db, logger)
	authService := authService.NewAuthService(authRepository, logger)
	authCtrl := authController.NewAuthController(authService)

	// Alat Dependencies
	alatRepository := alatRepo.NewAlatRepository(db, logger)
	alatService := alatService.NewAlatService(alatRepository, logger)
	alatCtrl := alatController.NewAlatController(alatService)

	// Raw Dependencies
	rawRepository := rawRepo.NewRawRepository(db, logger)
	rawService := rawService.NewRawService(rawRepository, logger)
	rawCtrl := rawController.NewRawController(rawService)

	// Validation Dependencies
	validRepository := validRepo.NewValidationRepository(db, logger)
	validService := validServ.NewValidationService(validRepository, logger)
	validCtrl := validCont.NewValidationController(validService, logger)

	// Maintenance Dependencies
	maintRepository := maintRepo.NewMaintenanceRepository(db, logger)
	maintService := maintService.NewMaintenanceService(maintRepository, logger)
	maintCtrl := maintController.NewMaintenanceController(maintService)

	r := gin.Default()

	// CORS Middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Basic Routes
	r.GET("/home", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "HTTP server listening",
			"success": true,
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
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

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.POST("/alat", alatCtrl.CreateAlat)
		apiGroup.GET("/alat/lantai/:lantai", alatCtrl.GetByLantai)
		apiGroup.GET("/alat/stasiun/:stasiun", alatCtrl.GetByStasiun)
		// raw
		apiGroup.POST("/raws/:alatId", rawCtrl.CreateByAlatId)
		apiGroup.GET("/raws/:alatId", rawCtrl.GetById)
		//validation
		apiGroup.POST("/validation/process", validCtrl.ValidateDataHandler)
		apiGroup.POST("/validation/process-batch", validCtrl.ValidateBatchHandler)
		apiGroup.GET("/validation/health", validCtrl.HealthCheckHandler)
		//maintenance
		apiGroup.POST("/maintenance/", middleware.JWTMiddleware(),  maintCtrl.Create)
		apiGroup.PUT("/maintenance/:id", middleware.JWTMiddleware(), maintCtrl.Update)
		apiGroup.DELETE("/maintenance/:id", middleware.JWTMiddleware(), maintCtrl.Delete)
		apiGroup.GET("/maintenance/:id", middleware.JWTMiddleware(),  maintCtrl.GetLogById)
		apiGroup.GET("/maintenance/lantai/:lantai", middleware.JWTMiddleware(), maintCtrl.GetLogByLantai)
		// admin
		apiGroup.POST("admin/register", authCtrl.Register)
		apiGroup.POST("admin/login", authCtrl.Login)

	}

	return r
}
