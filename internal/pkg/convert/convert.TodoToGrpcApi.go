package convert

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/chapar-rest/mock-server/internal/gen/mockv1"
	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

// TodoToGrpcApi converts a todo to a gRPC api todo.
func TodoToGrpcApi(t *model.Todo) *mockv1.Todo {
	return &mockv1.Todo{
		Id:          string(t.Id),
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    TodoPriorityToGrpcApi(t.Priority),
		Tags:        t.Tags,
		DueAt:       timeToGrpcApi(t.DueAt),
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}
}

// TodoPriorityToGrpcApi converts a todo priority to its gRPC enum.
func TodoPriorityToGrpcApi(p model.TodoPriority) mockv1.TodoPriority {
	switch p {
	case model.TodoPriorityLow:
		return mockv1.TodoPriority_TODO_PRIORITY_LOW
	case model.TodoPriorityMedium:
		return mockv1.TodoPriority_TODO_PRIORITY_MEDIUM
	case model.TodoPriorityHigh:
		return mockv1.TodoPriority_TODO_PRIORITY_HIGH
	}
	return mockv1.TodoPriority_TODO_PRIORITY_UNSPECIFIED
}

// TodoPriorityFromGrpcApi converts a gRPC priority enum to a todo priority.
// Unspecified maps to the empty priority, which the service treats as its
// default; unknown values map to an invalid priority the service rejects.
func TodoPriorityFromGrpcApi(p mockv1.TodoPriority) model.TodoPriority {
	switch p {
	case mockv1.TodoPriority_TODO_PRIORITY_UNSPECIFIED:
		return ""
	case mockv1.TodoPriority_TODO_PRIORITY_LOW:
		return model.TodoPriorityLow
	case mockv1.TodoPriority_TODO_PRIORITY_MEDIUM:
		return model.TodoPriorityMedium
	case mockv1.TodoPriority_TODO_PRIORITY_HIGH:
		return model.TodoPriorityHigh
	}
	return model.TodoPriority(p.String())
}

// TimeFromGrpcApi converts an optional protobuf timestamp.
func TimeFromGrpcApi(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}
	t := ts.AsTime().UTC()
	return &t
}

func timeToGrpcApi(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
