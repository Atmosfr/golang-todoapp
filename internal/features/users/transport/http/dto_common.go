package users_transport_http

import "github.com/Atmosfr/golang-todoapp/internal/core/domain"

type UserDTOResponse struct {
	ID          int     `json:"id"           example:"1"`
	Version     int     `json:"version"      example:"1"`
	FullName    string  `json:"full_name"    example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" example:"+79998887766"`
}

func userDTOFromDomain(userDomain domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          userDomain.ID,
		Version:     userDomain.Version,
		FullName:    userDomain.FullName,
		PhoneNumber: userDomain.PhoneNumber,
	}
}

func usersDTOFromDomains(usersDomains []domain.User) []UserDTOResponse {
	usersDTOResponse := make([]UserDTOResponse, len(usersDomains))
	for i, user := range usersDomains {
		usersDTOResponse[i] = userDTOFromDomain(user)
	}
	return usersDTOResponse
}
