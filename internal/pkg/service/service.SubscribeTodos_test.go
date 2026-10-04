package service

import (
	"testing"

	"github.com/chapar-rest/mock-server/internal/pkg/model"
)

func TestTodoHubDropsSlowSubscriber(t *testing.T) {
	h := newTodoHub()
	slow, cancelSlow := h.subscribe("s")
	defer cancelSlow()
	other, cancelOther := h.subscribe("other")
	defer cancelOther()

	//nolint:exhaustruct // only the type matters here.
	ev := model.TodoEvent{Type: model.TodoEventReset}
	for range todoEventBuffer + 1 {
		h.publish("s", ev)
	}

	got := 0
	for range slow {
		got++
	}
	if got != todoEventBuffer {
		t.Fatalf("slow subscriber got %d events before being dropped, want %d", got, todoEventBuffer)
	}
	select {
	case e := <-other:
		t.Fatalf("other session got %+v", e)
	default:
	}

	// Cancelling after being dropped is harmless.
	cancelSlow()
}
