package auth

import (
	"context"

	ssov1 "github.com/dm1tl/protos/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

func RegisterAuthServer(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	loginRequest *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{Token: "DSFDSF"}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	registerRequest *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	panic("implement me!")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	isAdminRequest *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	panic("implement me!")
}
