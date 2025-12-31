package main

import (
	"fmt"
	"log"

	"app/config"
	_ "app/docs" // Import generated docs
	"app/internal/adapters/broker/kafka"
	"app/internal/adapters/handler/http"
	"app/internal/adapters/storage/postgres"
	"app/internal/adapters/storage/redis"
	"app/internal/core/services"
	"app/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title AI Camera Backend API
// @version 1.0
// @description API Server for AI Camera System
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.InitLogger("development") // Change to "production" based on env
	defer logger.Log.Sync()
	logger.Info("Starting AI Camera API Server...")

	// 3. Init Database
	db, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		// don't die here if you want to support running without db for limited features,
		// but usually we fail fast.
		return
	}
	defer db.Close()

	// 4. Init Redis
	rdb, err := redis.NewRedisClient(cfg.Redis)
	if err != nil {
		logger.Error("Failed to connect to redis", zap.Error(err))
		return
	}
	defer rdb.Close()

	// 5. Init Kafka
	producer := kafka.NewProducer(cfg.Kafka)
	defer producer.Close()

	// 6. Init Router
	r := gin.Default()

	// --- CORS CONFIGURATION ---
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"db":      "connected",
			"redis":   "connected",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- WIRING DEPENDENCIES ---

	// Repositories
	cameraRepo := postgres.NewCameraRepository(db)
	userRepo := postgres.NewUserRepository(db)
	zoneRepo := postgres.NewZoneRepository(db)
	identityRepo := postgres.NewIdentityRepository(db)

	// Services
	cameraService := services.NewCameraService(cameraRepo)
	authService := services.NewAuthService(userRepo)
	zoneService := services.NewZoneService(zoneRepo)
	identityService := services.NewIdentityService(identityRepo)

	// Handlers
	cameraHandler := http.NewCameraHandler(cameraService)
	authHandler := http.NewAuthHandler(authService)
	zoneHandler := http.NewZoneHandler(zoneService)
	identityHandler := http.NewIdentityHandler(identityService)

	// --- ROUTES ---
	apiV1 := r.Group("/api/v1")
	{
		// Auth Routes (Public)
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected Routes
		protected := apiV1.Group("/")
		protected.Use(http.AuthMiddleware())
		{
			// Zones
			zones := protected.Group("/zones")
			{
				zones.POST("", zoneHandler.CreateZone)
				zones.GET("", zoneHandler.ListZones)
				zones.GET("/:id", zoneHandler.GetZone)
				zones.PUT("/:id", zoneHandler.UpdateZone)
				zones.DELETE("/:id", zoneHandler.DeleteZone)
			}

			// Cameras
			cameras := protected.Group("/cameras")
			{
				cameras.POST("", cameraHandler.CreateCamera)
				cameras.GET("", cameraHandler.ListCameras)
				cameras.GET("/:id", cameraHandler.GetCamera)
				cameras.PUT("/:id", cameraHandler.UpdateCamera)
				cameras.DELETE("/:id", cameraHandler.DeleteCamera)
			}

			// Identities
			identities := protected.Group("/identities")
			{
				identities.POST("", identityHandler.CreateIdentity)
				identities.GET("", identityHandler.ListIdentities)
				identities.GET("/:id", identityHandler.GetIdentity)
				identities.PUT("/:id", identityHandler.UpdateIdentity)
				identities.PATCH("/:id/status", identityHandler.UpdateStatus)
				identities.DELETE("/:id", identityHandler.DeleteIdentity)
			}
		}
	}

	// 7. Start Server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server listening", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
