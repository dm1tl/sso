package authserv

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	tokenTTL time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New returns new instance of Auth service
func New(
	log *slog.Logger,
	usrSvr UserSaver,
	usrPrvdr UserProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:      log,
		usrSvr:   usrSvr,
		usrPrvdr: usrPrvdr,
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
}

// Login check if user exists in db, if exists, returns token
// if doesn't exist or user input incorrect password - error
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string) (string, error) {
	const op = "internal.services.auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)
	if err := godotenv.Load(); err != nil {
		log.Error("couldn't load env variables %w", slogerr.Err(err))
	}
	log.Info("trying to login user")
	user, err := a.usrPrvdr.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrNotExists) {
			log.Warn("user not found", slogerr.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("couldn't get user with such id", slogerr.Err(err))
		return "", fmt.Errorf("op: %s", op)
	}

	log.Info("getting password from db")
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Info("invalid credentionals", slogerr.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in succesfully")
	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken check token's actuality, if token is actual - return id,
// else - return 0 and error
func (a *Auth) ValidateToken(ctx context.Context, token string) (UserId int64, err error) {
	const op = "internal.services.auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("validating token")
	id, err := jwt.ParseToken(token)
	if err != nil {
		log.Info("invalid token", slogerr.Err(err))
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}
	log.Info("token validated")
	return id, nil
}

// RegisterNewUser input new user in db and returns token
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
		return 0, fmt.Errorf("op: %s", op)
	}

	id, err := a.usrSvr.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			log.Warn("such user is already exists", slogerr.Err(err))
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}
		log.Error("couldn't save password in db", slogerr.Err(err))
		return 0, fmt.Errorf("op: %s", op)
	}
	log.Info("user registered")
	return id, nil
}
