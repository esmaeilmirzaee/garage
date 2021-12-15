package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Create inserts a new user into the database
func Create(ctx context.Context, db *sqlx.DB, ns NewUser, now time.Time) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(ns.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "user creation; failed to hash the password")
	}

	user := User{
		ID:        uuid.New().String(),
		Name:      ns.Name,
		Email:     ns.Email,
		Roles:     ns.Roles,
		Password:  hash,
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
	}

	const q = `INSERT INTO users (user_id, name, email, roles, password, created_at, 
updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING;`

	if _, err := db.ExecContext(ctx, q, user.ID, user.Name, user.Email, user.Roles, user.Password, user.CreatedAt,
		user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
