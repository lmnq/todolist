package service

import (
	"context"

	"github.com/lmnq/todolist/internal/entity"
)

type Service interface {
	CreateTask(ctx context.Context, task *entity.Task) (string, error)
	UpdateTask(ctx context.Context, id string, task *entity.Task) error
	DeleteTask(ctx context.Context, id string) error
	SetTaskStatusDone(ctx context.Context, id string) error
	GetTaskListByStatus(ctx context.Context, status string) ([]*entity.Task, error)
}
