package main

import (
	"github.com/kavya/content-engine/internal/config"
	"github.com/kavya/content-engine/internal/logger"
	"github.com/kavya/content-engine/internal/handlers"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.ErrorLog.Println("Failed to load configuration: %v", err)
	}

	logger.Init()

	logger.InfoLog.Println("Loaded Server port:", cfg.Server.Port)
	logger.InfoLog.Println("Loaded Server environment:", cfg.Server.Environment)

	handlers.SetupRoutes(cfg.Server.Port)
}
