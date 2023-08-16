package service

import (
	"context"
	"fmt"
	"time"

	"github.com/lmnq/todolist/internal/entity"
	"github.com/lmnq/todolist/internal/repo"
)

type ToDoListService struct {
	r repo.Repository
}

func New(r repo.Repository) *ToDoListService {
	return &ToDoListService{r: r}
}

func (s *ToDoListService) CreateTask(ctx context.Context, task *entity.Task) (string, error) {
	if err := task.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate task: %w", err)
	}

	id, err := s.r.CreateTask(ctx, task)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *ToDoListService) UpdateTask(ctx context.Context, id string, task *entity.Task) error {
	if err := task.Validate(); err != nil {
		return fmt.Errorf("failed to validate task: %w", err)
	}

	err := s.r.UpdateTask(ctx, id, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *ToDoListService) DeleteTask(ctx context.Context, id string) error {
	err := s.r.DeleteTask(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ToDoListService) SetTaskStatusDone(ctx context.Context, id string) error {
	err := s.r.SetTaskStatusDone(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ToDoListService) GetTaskListByStatus(ctx context.Context, status string) ([]*entity.Task, error) {
	tasks, err := s.r.GetTaskListByStatus(ctx, status)
	if err != nil {
		return []*entity.Task{}, err
	}

	for _, task := range tasks {
		weekday := task.ActiveAt.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			task.Title = "ВЫХОДНОЙ - " + task.Title
		}
	}

	return tasks, nil
}
