package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "internal.storage.sqlite.new()"
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) DeleteUser(ctx context.Context,
	id int64) (err error) {
	const op = "internal.storage.sqlite.DeleteUser()"
	stmt, err := s.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsCount, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsCount != 1 {
		return errors.New("user hasn't deleted")
	}
	return nil
}

func (s *Storage) SaveUser(ctx context.Context,
	email string,
	passwordHash []byte) (uid int64, err error) {
	const op = "internal.storage.sqlite.SaveUSer()"
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) values(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, email, passwordHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrAlreadyExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context,
	email string) (models.User, error) {
	const op = "internal.storage.sqlite.User()"
	var user models.User
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return user, fmt.Errorf("%s: %w", op, err)
	}
	row := stmt.QueryRowContext(ctx, email)
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrNotExists)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}
