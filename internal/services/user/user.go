package userserv

import (
	"context"
	"log/slog"
	"time"
)

type User struct {
	log      *slog.Logger
	usrDltr  userDeleter
	tokenTTL time.Duration
}

// New returns new instance of User service
func New(
	log *slog.Logger,
	usrDltr userDeleter,
	tokenTTL time.Duration,
) *User {
	return &User{
		log:      log,
		usrDltr:  usrDltr,
		tokenTTL: tokenTTL,
	}
}

type userDeleter interface {
	DeleteUser(ctx context.Context,
		id int64) (err error)
}

func (a *User) DeleteUser(ctx context.Context,
	id int64) (err error) {
	const op = "internal.services.user.DeletUser"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("trying to delete user")
	err = a.usrDltr.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
