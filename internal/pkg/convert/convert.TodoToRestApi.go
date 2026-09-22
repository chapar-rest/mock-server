package convert

import (
	"github.com/google/uuid"

	"github.com/chapar-rest/mock-server/internal/gen/restapi"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// TodoToRestApi converts a todo to a REST api todo.
func TodoToRestApi(t *model.Todo) restapi.Todo {
	return restapi.Todo{
		Id:          uuid.MustParse(string(t.Id)),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    restapi.TodoPriority(t.Priority),
		Tags:        t.Tags,
		DueAt:       t.DueAt,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}
