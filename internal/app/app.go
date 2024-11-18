package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	authserv "sso/internal/services/auth"
	userserv "sso/internal/services/user"
	"sso/internal/storage/sqlite"
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
	storage, err := sqlite.New(storagepath)
	if err != nil {
		panic(err)
	}
	authService := authserv.New(log, storage, storage, tokenTTL)
	userService := userserv.New(log, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, userService, grpcport)
	return &App{
		GRPCServer: grpcApp,
	}
}
