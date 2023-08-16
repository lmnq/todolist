package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lmnq/todolist/internal/service"
	"golang.org/x/exp/slog"
)

func NewRouter(router *chi.Mux, sl *slog.Logger, s service.Service) {
	// middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	// controller
	tc := newToDoListController(sl, s)

	// routes
	router.Route("/api/todo-list/tasks", func(r chi.Router) {
		r.Post("/", tc.createTask)
		r.Put("/{id}", tc.updateTask)
		r.Delete("/{id}", tc.deleteTask)
		r.Put("/{id}/done", tc.setTaskStatusDone)
		r.Get("/", tc.getTaskListByStatus)
	})
}
