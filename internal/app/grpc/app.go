package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	authgrpc "sso/internal/grpc/auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCserver *grpc.Server
	port       string
}

func New(logger *slog.Logger, port string) *App {
	grpcServer := grpc.NewServer()
	authgrpc.RegisterAuthServer(grpcServer)
	return &App{
		log:        logger,
		gRPCserver: grpcServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	op := "internal.app.grpc.New()"
	log := a.log.With(slog.String("op", op),
		slog.String("port", a.port))
	lis, err := net.Listen("tcp", ":"+a.port)
	if err != nil {
		return fmt.Errorf("couldn't listen tcp: error %v, op %v", err, op)
	}
	log.Info("grpc server is running", slog.String("grpc address", lis.Addr().String()))
	if err := a.gRPCserver.Serve(lis); err != nil {
		return fmt.Errorf("couldn't listen tcp: error %v, op %v", err, op)
	}

	return nil
}

func (a *App) Stop() {
	op := "internal.app.grpc.Stop()"
	log := a.log.With(slog.String("op", op),
		slog.String("port", a.port))
	log.Info("stopping grpc server")
	a.gRPCserver.GracefulStop()
}
