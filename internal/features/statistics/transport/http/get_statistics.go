package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/response"
)

type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"                 example:"10"`
	TasksCompleted             int      `json:"tasks_completed"               example:"5"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"          example:"50"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"20m30s"`
}

// GetStatistics godoc
// @Summary Get statistics
// @Description Get statistics with optional filtering by user ID and date range.
// @Tags statistics
// @Produce json
// @Param user_id query int false "Author ID"
// @Param from query string false "Start date (inclusive), format YYYY-MM-DD"
// @Param to query string false "End date (exclusive), format YYYY-MM-DD"
// @Success 200 {object} GetStatisticsResponse "Statistics retrieved successfully"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userID, from, to, err := getStatisticsQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/from/to query params")
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := statisticsDTOFromDomain(statisticsDomain)

	responseHandler.JSONResponse(response, http.StatusOK)
}

func statisticsDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		avgTimeStr := statistics.TasksAverageCompletionTime.String()
		avgTime = &avgTimeStr
	}
	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TasksCompletedRate:         statistics.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getStatisticsQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		userIDQueryParam = "user_id"
		fromQueryParam   = "from"
		toQueryParam     = "to"
	)
	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}
	from, err := core_http_request.GetDateQueryParam(r, fromQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'from' date query param: %w", err)
	}
	to, err := core_http_request.GetDateQueryParam(r, toQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'to' date query param: %w", err)
	}

	return userID, from, to, nil
}
