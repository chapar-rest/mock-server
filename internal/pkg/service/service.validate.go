package service

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

const (
	maxTitleLength       = 200
	maxDescriptionLength = 2000
	maxTags              = 10
	maxTagLength         = 32
)

// validateTodo checks the user-editable fields of a todo.
func validateTodo(t *model.Todo) error {
	if strings.TrimSpace(t.Title) == "" {
		return fmt.Errorf("%w: title is required", errx.ErrInvalidArgument)
	}
	if utf8.RuneCountInString(t.Title) > maxTitleLength {
		return fmt.Errorf("%w: title must be at most %d characters", errx.ErrInvalidArgument, maxTitleLength)
	}
	if utf8.RuneCountInString(t.Description) > maxDescriptionLength {
		return fmt.Errorf("%w: description must be at most %d characters", errx.ErrInvalidArgument, maxDescriptionLength)
	}
	if _, err := model.ParseTodoPriority(string(t.Priority)); err != nil {
		return err
	}
	if len(t.Tags) > maxTags {
		return fmt.Errorf("%w: at most %d tags are allowed", errx.ErrInvalidArgument, maxTags)
	}
	for _, tag := range t.Tags {
		if tag == "" || utf8.RuneCountInString(tag) > maxTagLength {
			return fmt.Errorf("%w: tags must be 1-%d characters", errx.ErrInvalidArgument, maxTagLength)
		}
	}
	return nil
}
