package users_service

import (
	"context"
	"fmt"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {

	user, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user from repository: %w", err)
	}

	return user, nil
}
