package main

import (
	"internet_shop/internal/app"
	"internet_shop/internal/config"
	"internet_shop/pkg/logging"
	"log"
)

func main() {
	log.Print("config initializing")
	cfg := config.GetConfig()

	log.Print("logger initializing")
	logger := logging.GetLogger(cfg.AppConfig.LogLevel)
	
	a, err := app.NewApp(logger, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Println("raning aplication")
	a.Run()
}