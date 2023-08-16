package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lmnq/todolist/internal/entity"
	"github.com/lmnq/todolist/internal/service"
	"github.com/lmnq/todolist/internal/utils"
	"golang.org/x/exp/slog"
)

type toDoListController struct {
	sl *slog.Logger
	s  service.Service
}

func newToDoListController(sl *slog.Logger, s service.Service) *toDoListController {
	return &toDoListController{sl, s}
}

type taskRequest struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
}

type createTaskResponse struct {
	ID string `json:"id"`
}

func (c *toDoListController) createTask(w http.ResponseWriter, r *http.Request) {
	var reqTask taskRequest

	err := json.NewDecoder(r.Body).Decode(&reqTask)
	if err != nil {
		c.sl.Error("failed to decode request body: %v", err)
		http.NotFound(w, r)
		return
	}

	date, err := utils.ParseTime(reqTask.ActiveAt)
	if err != nil {
		c.sl.Error("failed to parse date: %v", err)
		http.NotFound(w, r)
		return
	}

	id, err := c.s.CreateTask(r.Context(), &entity.Task{
		Title:    reqTask.Title,
		ActiveAt: date,
	})
	if err != nil {
		c.sl.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createTaskResponse{ID: id})
}

func (c *toDoListController) updateTask(w http.ResponseWriter, r *http.Request) {
	var reqTask taskRequest

	err := json.NewDecoder(r.Body).Decode(&reqTask)
	if err != nil {
		c.sl.Error("failed to decode request body: %v", err)
		http.NotFound(w, r)
		return
	}

	date, err := utils.ParseTime(reqTask.ActiveAt)
	if err != nil {
		c.sl.Error("failed to parse date: %v", err)
		http.NotFound(w, r)
		return
	}

	id := chi.URLParam(r, "id")

	err = c.s.UpdateTask(r.Context(), id, &entity.Task{
		Title:    reqTask.Title,
		ActiveAt: date,
	})
	if err != nil {
		c.sl.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *toDoListController) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.s.DeleteTask(r.Context(), id)
	if err != nil {
		c.sl.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *toDoListController) setTaskStatusDone(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := c.s.SetTaskStatusDone(r.Context(), id)
	if err != nil {
		c.sl.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type taskListResponse struct {
	Title    string `json:"title"`
	ActiveAt string `json:"activeAt"`
}

func (c *toDoListController) getTaskListByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	tasks, err := c.s.GetTaskListByStatus(r.Context(), status)
	if err != nil {
		c.sl.Error(err.Error())
		http.NotFound(w, r)
		return
	}

	var result []taskListResponse
	for _, task := range tasks {
		fmt.Println(task.Title)
		result = append(result, taskListResponse{
			Title:    task.Title,
			ActiveAt: utils.FormatTime(task.ActiveAt),
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
