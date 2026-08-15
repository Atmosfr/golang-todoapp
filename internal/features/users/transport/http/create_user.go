package users_transport_http

import (
	"net/http"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
	core_errors "github.com/Atmosfr/golang-todoapp/internal/core/errors"
	core_logger "github.com/Atmosfr/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Atmosfr/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name"    validate:"required,min=3,max=100"               example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79995551423"`
}

type CreateUserResponse UserDTOResponse

// CreateUser godoc
// @Summary Create a user
// @Description Create a new user in the system.
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "CreateUser request payload"
// @Success 201 {object} CreateUserResponse "User created successfully"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router /users [post]
func (h *UsersHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := core_logger.FromContext(ctx)
	logger.Debug("invoke CreateUser handler")
	responseHandler := core_http_response.NewHTTPResponseHandler(logger, w)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)
	userResponseDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userResponseDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
