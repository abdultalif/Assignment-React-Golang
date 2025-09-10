package app

import (
	"backend-service/config"
	"backend-service/internal/logger"
	"log"

	"github.com/gin-gonic/gin"
)

func RunServer() {
	cfg := config.NewConfig()

	_, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-2] failed to connect postgres: %v", err)
		return
	}

	logger.InitLogger()

	r := gin.Default()

	r.Run(cfg.App.AppPort)
}
