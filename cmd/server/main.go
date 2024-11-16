package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/innoglobe/xmgo/docs"
	"github.com/innoglobe/xmgo/internal/app"
	"github.com/innoglobe/xmgo/internal/config"
	postgresrepository "github.com/innoglobe/xmgo/internal/infrastructure/db/postgres"
	"github.com/innoglobe/xmgo/internal/infrastructure/server"
	"github.com/innoglobe/xmgo/internal/infrastructure/server/handler"
	eventservice "github.com/innoglobe/xmgo/internal/service"
	"github.com/innoglobe/xmgo/internal/usecase"
	"github.com/innoglobe/xmgo/pkg/logger"
	"github.com/innoglobe/xmgo/pkg/migrations"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// @title XMGO API
// @version 1.0
// @description This is a sample server for XMGO.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	// Config file
	configFile := flag.String("config", "/config/config.yaml", "Path to config file")
	flag.Parse()
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize kafka producer
	kafkaProducer := eventservice.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	// Initialize logger
	log := logger.NewLogger()

	// Initialize db conn
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error(fmt.Sprint(err.Error()))
	}

	// Run database migrations
	log.Info("Starting migrations...")
	if err := migrations.Migrate(dsn, "./migrations"); err != nil {
		log.Error(fmt.Sprintf("Failed to run migrations: %v", err))
	} else {
		log.Info("Migrations applied successfully.")
	}

	// Initialize repository and usecase
	companyRepo := postgresrepository.NewPostgresRepository(db)
	companyUsecase := usecase.NewCompanyUsecase(companyRepo, kafkaProducer)

	// Initialize handlers
	companyHandler := handler.NewCompanyHandler(companyUsecase)
	authHandler := handler.NewAuthHandler(cfg.JWT.Secret)

	// Switch gin to release mode if needed
	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router with handler
	r := server.NewRouter(companyHandler, authHandler)
	router := r.RegisterRoutes(cfg.JWT.Secret)

	// Configure CORS
	//router.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"*"},
	//	AllowMethods:     []string{"POST", "PATCH", "DELETE"},
	//	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}))

	// Add security headers
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Next()
	})

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize the application
	application, err := app.NewApp(cfg, companyUsecase, log, router, kafkaProducer)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to initialize app: %v", err))
	}

	// Run the application
	if err := application.Run(); err != nil {
		log.Error(fmt.Sprintf("Application failed: %v", err))
	}
}
