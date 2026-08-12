package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
	core_errors "github.com/Atmosfr/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || to.Equal(*from) {
			return domain.Statistics{}, fmt.Errorf(
				"`to` must be after `from`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.statisticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("failed to get statistics from repository: %w", err)
	}

	statistics := calcStatistics(tasks)

	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	tasksCreated := len(tasks)
	if tasksCreated == 0 {
		return domain.Statistics{}
	}

	tasksCompleted := 0
	var totalCompletionDuration time.Duration
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++

			taskCompletion := task.CompletionDuration()
			if taskCompletion != nil {
				totalCompletionDuration += *taskCompletion
			}
		}
	}

	rate := float64(tasksCompleted) / float64(tasksCreated) * 100
	tasksCompletedRate := &rate

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		tasksCompletedRate,
		tasksAverageCompletionTime,
	)
}
