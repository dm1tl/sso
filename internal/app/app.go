package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcport string,
	storagepath string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcport)
	return &App{
		GRPCServer: grpcApp,
	}
}
