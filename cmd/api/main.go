package main

import (
	"fmt"
	"log"
	"os"

	"app/config"
	_ "app/docs" // Import generated docs
	"app/internal/adapters/broker/kafka"
	"app/internal/adapters/handler/http"
	localstorage "app/internal/adapters/storage/local"
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

// @title           AI Camera API
// @version         1.0
// @description     This is the API server for AI Camera System.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

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
	logger.InitLogger("development")
	defer logger.Log.Sync()
	logger.Info("Starting AI Camera API Server...")

	// 3. Init Database
	db, err := postgres.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve Static Files (Uploads)
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}
	r.Static("/uploads", uploadDir)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong", "db": "connected", "redis": "connected"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- WIRING DEPENDENCIES ---

	// Repositories & Adapters
	cameraRepo := postgres.NewCameraRepository(db)
	userRepo := postgres.NewUserRepository(db)
	zoneRepo := postgres.NewZoneRepository(db)
	identityRepo := postgres.NewIdentityRepository(db)
	faceRepo := postgres.NewIdentityFaceRepository(db)
	aiRepo := postgres.NewAIRepository(db)
	roleRepo := postgres.NewRoleRepository(db)
	analyticsRepo := postgres.NewAnalyticsRepository(db)
	auditRepo := postgres.NewAuditRepository(db)
	permRepo := postgres.NewPermissionRepository(db)

	// Host for static files
	baseURL := fmt.Sprintf("http://localhost:%d/uploads", cfg.Server.Port)
	fileStorage := localstorage.NewLocalStorage(uploadDir, baseURL)

	// Services
	cameraService := services.NewCameraService(cameraRepo)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo) // Added UserService
	zoneService := services.NewZoneService(zoneRepo)
	identityService := services.NewIdentityService(identityRepo, faceRepo)
	aiService := services.NewAIService(aiRepo)
	roleService := services.NewRoleService(roleRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo)
	auditService := services.NewAuditService(auditRepo)
	permService := services.NewPermissionService(permRepo)
	mediaService := services.NewMediaService(fileStorage)

	// Handlers
	cameraHandler := http.NewCameraHandler(cameraService)
	authHandler := http.NewAuthHandler(authService)
	userHandler := http.NewUserHandler(userService) // Added UserHandler
	zoneHandler := http.NewZoneHandler(zoneService)
	identityHandler := http.NewIdentityHandler(identityService)
	aiHandler := http.NewAIHandler(aiService)
	roleHandler := http.NewRoleHandler(roleService)
	analyticsHandler := http.NewAnalyticsHandler(analyticsService)
	auditHandler := http.NewAuditHandler(auditService)
	permHandler := http.NewPermissionHandler(permService)
	mediaHandler := http.NewMediaHandler(mediaService)

	// --- ROUTES ---
	apiV1 := r.Group("/api/v1")
	{
		auth := apiV1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := apiV1.Group("/")
		protected.Use(http.AuthMiddleware())
		{
			// Media Upload
			protected.POST("/media/upload", mediaHandler.UploadImage)

			// Dashboard & AI
			protected.GET("/stats/dashboard", aiHandler.GetDashboardStats)
			protected.GET("/ai-configs/camera/:cameraId", aiHandler.GetConfig)
			protected.POST("/ai-configs", aiHandler.UpdateConfig)
			protected.PUT("/ai-configs/:id", aiHandler.UpdateConfig)
			protected.GET("/events", aiHandler.ListEvents)
			protected.PATCH("/events/:id", aiHandler.UpdateEventStatus)

			// Analytics & Attendance
			analytics := protected.Group("")
			{
				analytics.GET("/recognition/logs", analyticsHandler.ListRecognitionLogs)
				analytics.GET("/attendance/records", analyticsHandler.ListAttendance)
				analytics.GET("/attendance/summary", analyticsHandler.GetSummary)
			}

			// System Logs
			protected.GET("/audit-logs", auditHandler.ListLogs)

			// Permissions (Data Scoping)
			protected.GET("/permissions/:userId", permHandler.GetPermissions)
			protected.POST("/permissions/:userId/cameras", permHandler.UpdateCameraPermissions)
			protected.POST("/permissions/:userId/zones", permHandler.UpdateZonePermissions)

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

			// Identities & Faces
			identities := protected.Group("/identities")
			{
				identities.POST("", identityHandler.CreateIdentity)
				identities.GET("", identityHandler.ListIdentities)
				identities.GET("/:id", identityHandler.GetIdentity)
				identities.PUT("/:id", identityHandler.UpdateIdentity)
				identities.PATCH("/:id/status", identityHandler.UpdateStatus)
				identities.DELETE("/:id", identityHandler.DeleteIdentity)

				identities.POST("/enroll-face", identityHandler.EnrollFace)
				identities.DELETE("/faces/:face_id", identityHandler.DeleteFace)
			}

			// Roles
			roles := protected.Group("/roles")
			{
				roles.POST("", roleHandler.CreateRole)
				roles.GET("", roleHandler.ListRoles)
				roles.GET("/:id", roleHandler.GetRole)
				roles.PUT("/:id", roleHandler.UpdateRole)
				roles.DELETE("/:id", roleHandler.DeleteRole)
			}

			// Users configuration
			users := protected.Group("/users")
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.ListUsers)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
				users.POST("/:user_id/reset-password", userHandler.ResetPassword)
			}
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("Server listening", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Error("Failed to run server", zap.Error(err))
	}
}
