package app

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lmnq/todolist/config"
	"github.com/lmnq/todolist/internal/controller"
	"github.com/lmnq/todolist/internal/repo"
	"github.com/lmnq/todolist/internal/service"
	"github.com/lmnq/todolist/pkg/mongodb"
	"golang.org/x/exp/slog"
)

const (
	timeout = 10 * time.Second
)

func Run(cfg *config.Config) {
	// создание логгера
	var sl *slog.Logger = slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	// подключение к БД
	mongoClient, err := mongodb.NewMongoClient(cfg.MongoDB.URI)
	if err != nil {
		sl.Error("mongodb error: %v", err)
		os.Exit(1)
	}
	defer mongoClient.Disconnect(context.Background())

	// создание индексов
	err = repo.CreateUniqueIndexes(mongoClient)
	if err != nil {
		sl.Error("failed to create unique indexes: %v", err)
		os.Exit(1)
	}

	// создание репозитория
	repository := repo.New(mongoClient)

	// создание сервиса
	toDoListService := service.New(repository)

	// создание роутера
	router := chi.NewRouter()

	// добавление маршрутов
	controller.NewRouter(router, sl, toDoListService)

	// HTTP сервер
	server := http.Server{
		Addr:         net.JoinHostPort("localhost", cfg.Port),
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		Handler:      router,
	}

	// запуск
	err = server.ListenAndServe()
	if err != nil {
		sl.Error("failed to start server: %v", err)
		os.Exit(1)
	}
}
