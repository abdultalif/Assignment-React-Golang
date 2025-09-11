package app

import (
	"backend-service/config"
	"backend-service/internal/adapter/handler"
	"backend-service/internal/adapter/repository"
	"backend-service/internal/adapter/router"
	"backend-service/internal/core/service"
	"backend-service/internal/logger"
	"backend-service/pkg/validator"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10/translations/en"
)

func RunServer() {
	r := gin.Default()
	cfg := config.NewConfig()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-2] failed to connect postgres: %v", err)
		return
	}

	customValidator := validator.NewValidator()
	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	r.Use(func(c *gin.Context) {
		c.Set("validator", customValidator)
		c.Next()
	})

	logger.InitLogger()

	productRepo := repository.NewProductRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)
	jobRepo := repository.NewJobRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)

	orderService := service.NewOrderService(orderRepo, productRepo)
	jobService := service.NewJobService(jobRepo, transactionRepo)

	orderHandler := handler.NewOrderHandler(orderService, customValidator)
	jobHandler := handler.NewJobHandler(jobService, customValidator)

	r = router.SetupRouter(orderHandler, jobHandler)

	log.Printf("Starting server on port %s", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
