package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks godoc
// @Summary Get a task list
// @Description Get a paginated list of tasks with optional filtering by author.
// @Tags tasks
// @Produce json
// @Param user_id query int false "Author ID"
// @Param limit query int false "Tasks per page"
// @Param offset query int false "Pagination offset"
// @Success 200 {object} GetTasksResponse "Task list retrieved successfully"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TasksHTTPHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'user_id', 'limit' or 'offset' query param",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	tasksDTOs := GetTasksResponse(taskDTOsFromDomains(tasksDomains))
	responseHandler.JSONResponse(tasksDTOs, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		user_idQueryParamKey = "user_id"
		limitQueryParamKey   = "limit"
		offsetQueryParamKey  = "offset"
	)

	user_id, err := core_http_request.GetIntQueryParam(r, user_idQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'user_id' query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return user_id, limit, offset, nil
}
