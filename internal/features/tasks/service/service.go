package tasks_service

import (
	"context"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
)

type TasksRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTasks(ctx context.Context, userID, limit, offset *int) ([]domain.Task, error)
	GetTask(ctx context.Context, taskID int) (domain.Task, error)
	DeleteTask(ctx context.Context, taskID int) error
	PatchTask(ctx context.Context, taskID int, task domain.Task) (domain.Task, error)
}

type TasksService struct {
	tasksRepository TasksRepository
}

func NewTasksService(tasksRepository TasksRepository) *TasksService {
	return &TasksService{
		tasksRepository: tasksRepository,
	}
}
