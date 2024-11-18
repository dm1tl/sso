package auth

import (
	"context"
	"errors"
	"sso/internal/lib/validation"
	authserv "sso/internal/services/auth"
	"sso/internal/storage"

	ssov1 "github.com/dm1tl/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string) (UserId int64, err error)
	ValidateToken(
		ctx context.Context,
		token string) (UserId int64, err error)
}

type User interface {
	DeleteUser(
		ctx context.Context,
		id int64) (err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	ssov1.UnimplementedUserServer
	auth Auth
	user User
}

func Register(gRPC *grpc.Server, auth Auth, user User) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth, user: user})
	ssov1.RegisterUserServer(gRPC, &serverAPI{auth: auth, user: user})
}

func (s *serverAPI) Delete(
	ctx context.Context,
	req *ssov1.DeleteRequest) (*ssov1.DeleteResponse, error) {
	err := s.user.DeleteUser(ctx, req.GetId())
	if err != nil {
		return &ssov1.DeleteResponse{ErrorMessage: "Unable to delete user"}, status.Error(codes.NotFound, "user not found")
	}
	return &ssov1.DeleteResponse{ErrorMessage: "success"}, nil
}
func (s *serverAPI) ValidateToken(
	ctx context.Context,
	req *ssov1.ValidateTokenRequest) (*ssov1.ValidateTokenResponse, error) {
	if err := validation.ValidateTokenData(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	id, err := s.auth.ValidateToken(ctx, req.GetToken())
	if err != nil {
		//TODO process error
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.ValidateTokenResponse{Id: id}, nil
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := validation.ValidateLoginData(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, authserv.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	if err := validation.ValidateRegisterData(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	userId, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &ssov1.RegisterResponse{UserId: userId}, nil
}
