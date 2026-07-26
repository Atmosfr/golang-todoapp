package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/Atmosfr/golang-todoapp/internal/core/errors"
)

type User struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(UninitializedID, UninitializedVersion, fullName, phoneNumber)
}

func NewUser(id, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func (u *User) Validate() error {
	fullNameLen := len([]rune(u.FullName))
	if fullNameLen < 3 || fullNameLen > 100 {
		return fmt.Errorf("invalid `FullName` len: %d: %w", fullNameLen, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneLen := len([]rune(*u.PhoneNumber))
		if phoneLen < 10 || phoneLen > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len: %d: %w", phoneLen, core_errors.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+\d+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("`FullName can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmpUser := *u

	if patch.FullName.Set {
		tmpUser.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmpUser.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmpUser.Validate(); err != nil {
		return fmt.Errorf("validate user after patch: %w", err)
	}

	*u = tmpUser
	return nil
}
