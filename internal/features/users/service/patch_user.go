package users_service

import (
	"context"
	"fmt"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	userID int,
	userPatch domain.UserPatch,
) (domain.User, error) {
	userDomain, err := s.usersRepository.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user by id: %w", err)
	}

	if err := userDomain.ApplyPatch(userPatch); err != nil {
		return domain.User{}, fmt.Errorf("apply user patch: %w", err)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, userID, userDomain)
	if err != nil {
		return domain.User{}, fmt.Errorf("patch user in repository: %w", err)
	}

	return patchedUser, nil
}
