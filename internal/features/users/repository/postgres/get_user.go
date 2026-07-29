package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
	core_errors "github.com/Atmosfr/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Atmosfr/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	userID int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, userID)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' not found in database: %w",
				userID,
				core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("get user scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.Version,
		userModel.FullName,
		userModel.PhoneNumber,
	)

	return userDomain, nil
}
