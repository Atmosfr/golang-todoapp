package web_fs_repository

import (
	"errors"
	"fmt"
	"os"

	core_errors "github.com/Atmosfr/golang-todoapp/internal/core/errors"
)

func (r *WebRepository) GetFile(path string) ([]byte, error) {
	file, err := os.ReadFile(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("file: %s: %w", path, core_errors.ErrNotFound)
		}

		return nil, fmt.Errorf("read file: %s: %w", path, err)
	}

	return file, nil
}
