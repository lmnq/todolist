package main

import (
	"log"

	"github.com/lmnq/todolist/config"
	"github.com/lmnq/todolist/internal/app"
)

func main() {
	// создание конфига
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	app.Run(cfg)
}
