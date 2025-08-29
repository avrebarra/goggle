package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	pb "github.com/avrebarra/goggle/proto/generated"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// mockable:true
type ToggleService servicetoggle.Service

type ConfigRuntime struct {
	Version       string        `validate:"required"`
	Port          int           `validate:"required"`
	ToggleService ToggleService `validate:"required"`
	StartedAt     time.Time     `validate:"required"`
}

type Runtime struct {
	Config ConfigRuntime
	Server *grpc.Server
}

func NewRuntime(cfg ConfigRuntime) (out *Runtime, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Wrap(err, "config validation failed")
		return
	}

	out = &Runtime{
		Config: cfg,
	}
	return
}

func (e *Runtime) Run() (err error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", e.Config.Port))
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}

	s := grpc.NewServer()
	handler := &Handler{ConfigRuntime: e.Config}

	pb.RegisterGoggleServiceServer(s, handler)
	reflection.Register(s) // Enable gRPC reflection for debugging

	e.Server = s

	slog.Info("gRPC server starting", "port", e.Config.Port)
	if err = s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed to serve gRPC")
	}

	return
}

func (e *Runtime) Start(ctx context.Context) <-chan bool {
	done := make(chan bool)
	go func() {
		defer close(done)
		if err := e.Run(); err != nil {
			slog.Error("gRPC server error", "error", err)
		}
		done <- true
	}()
	return done
}

func (e *Runtime) Stop() {
	if e.Server != nil {
		e.Server.GracefulStop()
	}
}
