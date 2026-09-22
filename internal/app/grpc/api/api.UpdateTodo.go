package api

import (
	"context"
	"fmt"
	"time"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/convert"
	"github.com/chapar-rest/mock-server/internal/pkg/errx"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
	"github.com/chapar-rest/mock-server/internal/pkg/service"
	"github.com/chapar-rest/mock-server/internal/pkg/utils"
)

// UpdateTodo implements mockv1.TodoServiceServer.
func (s *TodoServer) UpdateTodo(ctx context.Context, req *mockv1.UpdateTodoRequest) (*mockv1.Todo, error) {
	sessionId, err := sessionIdFromContext(ctx)
	if err != nil {
		return nil, err
	}

	in := req.GetTodo()
	if in == nil {
		return nil, fmt.Errorf("%w: todo is required", errx.ErrInvalidArgument)
	}

	todoId, err := model.ParseTodoId(in.GetId())
	if err != nil {
		return nil, err
	}

	var todo *model.Todo
	if paths := req.GetUpdateMask().GetPaths(); len(paths) == 0 {
		todo, err = s.service.ReplaceTodo(ctx, &service.ReplaceTodoRequest{
			SessionId:   sessionId,
			Id:          todoId,
			Title:       in.GetTitle(),
			Description: in.GetDescription(),
			Completed:   in.GetCompleted(),
			Priority:    convert.TodoPriorityFromGrpcApi(in.GetPriority()),
			Tags:        in.GetTags(),
			DueAt:       convert.TimeFromGrpcApi(in.GetDueAt()),
		})
	} else {
		var update *service.UpdateTodoRequest
		update, err = maskedUpdate(sessionId, todoId, in, paths)
		if err != nil {
			return nil, err
		}
		todo, err = s.service.UpdateTodo(ctx, update)
	}
	if err != nil {
		return nil, err
	}

	return convert.TodoToGrpcApi(todo), nil
}

// maskedUpdate copies the fields named by the update mask into a partial
// update.
func maskedUpdate(sessionId model.SessionId, todoId model.TodoId, in *mockv1.Todo, paths []string) (*service.UpdateTodoRequest, error) {
	update := &service.UpdateTodoRequest{
		SessionId:   sessionId,
		Id:          todoId,
		Title:       utils.None[string](),
		Description: utils.None[string](),
		Completed:   utils.None[bool](),
		Priority:    utils.None[model.TodoPriority](),
		Tags:        utils.None[[]string](),
		DueAt:       utils.None[*time.Time](),
	}

	for _, path := range paths {
		switch path {
		case "title":
			update.Title = utils.Some(in.GetTitle())
		case "description":
			update.Description = utils.Some(in.GetDescription())
		case "completed":
			update.Completed = utils.Some(in.GetCompleted())
		case "priority":
			update.Priority = utils.Some(convert.TodoPriorityFromGrpcApi(in.GetPriority()))
		case "tags":
			update.Tags = utils.Some(in.GetTags())
		case "due_at":
			update.DueAt = utils.Some(convert.TimeFromGrpcApi(in.GetDueAt()))
		default:
			return nil, fmt.Errorf("%w: update_mask path %q is not updatable", errx.ErrInvalidArgument, path)
		}
	}
	return update, nil
}
