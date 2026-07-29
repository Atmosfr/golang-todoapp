package users_postgres_repository

import (
	"github.com/Atmosfr/golang-todoapp/internal/core/domain"
)

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromModels(
	models []UserModel,
) []domain.User {
	domains := make([]domain.User, len(models))
	for i, model := range models {
		domains[i] = domain.NewUser(
			model.ID,
			model.Version,
			model.FullName,
			model.PhoneNumber,
		)
	}
	return domains
}
