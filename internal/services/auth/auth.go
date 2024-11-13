package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"sso/internal/domain/models"
	"sso/internal/lib/jwt"
	"sso/internal/lib/logger/slogerr"
	"sso/internal/storage"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log      *slog.Logger
	usrSvr   UserSaver
	usrPrvdr UserProvider
	appPrvdr AppProvider
	tokenTTL time.Duration
}

var (
	ErrInvalidCredentials = "invalid credentials"
)

// New returns new instance of Auth service
func New(
	log *slog.Logger,
	usrSvr UserSaver,
	usrPrvdr UserProvider,
	appPrvdr AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:      log,
		usrSvr:   usrSvr,
		usrPrvdr: usrPrvdr,
		appPrvdr: appPrvdr,
		tokenTTL: tokenTTL,
	}
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passwordHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context,
		email string) (models.User, error)
	IsAdmin(ctx context.Context,
		userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context,
		appId int64) (models.App, error)
}

// Login checks if user exists in db, if exists, returns token
// if doesn't exist or user input incorrect password - error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int64,
) (string, error) {
	const op = "internal.services.auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)
	if err := godotenv.Load(); err != nil {
		log.Error("couldn't load env variables %v", slogerr.Err(err))
	}
	log.Info("trying to login user")
	user, err := a.usrPrvdr.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			log.Warn("user not found", slogerr.Err(err))
			return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
		}
		log.Error("couldn't get user with such id", slogerr.Err(err))
		return "", fmt.Errorf("op: %v", op)
	}

	log.Info("getting password from db")
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Info("invalid credentionals", slogerr.Err(err))
		return "", fmt.Errorf("%s: %v", op, ErrInvalidCredentials)
	}

	log.Info("user logged in succesfully")
	secret := a.fetchJWTSecret()
	token, err := jwt.NewToken(user, secret, a.tokenTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Register input new user in db and returns token
// if fields aren't validated, or email is not unique returns error
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string) (UserId int64, err error) {
	const op = "internal.services.auth.Register"

	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("couldn't hash password", slogerr.Err(err))
		return 0, fmt.Errorf("op: %v", op)
	}

	id, err := a.usrSvr.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			log.Warn("such user is already exists", slogerr.Err(err))
			return 0, fmt.Errorf("%s: %v", op, storage.ErrAlreadyExists)
		}
		log.Error("couldn't save password in db", slogerr.Err(err))
		return 0, fmt.Errorf("op: %v", op)
	}
	log.Info("user registered")
	return id, nil
}

// Register checks if user with concrete id has user access or not
// if yes - returns true, else false
// if UserId is incorrect or doesn't exist in db returns error
func (a *Auth) IsAdmin(
	ctx context.Context,
	UserId int64) (bool, error) {
	panic("implement me")
}

func (a *Auth) fetchJWTSecret() string {
	const op = "internal.lib.jwt"
	log := a.log.With(
		slog.String("op", op),
	)
	if err := godotenv.Load(); err != nil {
		log.Error("couldn't load env variables %v", slogerr.Err(err))
	}
	res := os.Getenv("JWT_SECRET")
	return res
}