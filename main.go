package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/chapar-rest/mock-server/grpc"
	"github.com/chapar-rest/mock-server/rest"
)

func main() {
	restApi := rest.New()
	grpcApi := grpc.New()

	ctx := ContextWithOsSignal()

	go func() {
		_ = restApi.Run()
	}()

	go func() {
		_ = grpcApi.Run()
	}()

	<-ctx.Done()
}

var DefaultExistSignals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}

// ContextWithOsSignal returns a context with by default is listening to
// SIGHUP, SIGINT, SIGTERM, SIGQUIT os signals to cancel
func ContextWithOsSignal(sig ...os.Signal) context.Context {
	if len(sig) == 0 {
		sig = DefaultExistSignals
	}

	s := make(chan os.Signal, 1)
	signal.Notify(s, sig...)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-s
		cancel()
	}()
	return ctx
}
