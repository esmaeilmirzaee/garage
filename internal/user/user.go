package user

import (
	"context"
	"database/sql"
	"github.com/esmaeilmirzaee/grage/internal/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrAuthenticationFailure occurs when a User attempts to authenticate
	// but anything goes wrong.
	ErrAuthenticationFailure = errors.New("Authentication failed")
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

// Authenticate finds a User by their email and verifies their password. On
// success, it returns a Claims value representing this User. The Claims can be
// used to generate a token for future authentication.
func Authenticate(ctx context.Context, db *sqlx.DB, now time.Time, email, password string) (auth.Claims, error) {
	const q = `SELECT user_id, name, email, password, roles, created_at, updated_at FROM users WHERE email = $1;`
	var u User
	if err := db.GetContext(ctx, &u, q, email); err != nil {
		// Normally we would return ErrNotFound in this scenario, but we do not want
		// to leak to an unauthenticated user which emails are in the system
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the saved one. Use the bcrypt
	// comparison function, so it is cryptographically secure.
	if err := bcrypt.CompareHashAndPassword(u.Password, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// If we are this far the request is valid. Create some claims for the User
	// and generate their token.
	claims := auth.NewClaims(u.ID, u.Roles, now, time.Hour)
	return claims, nil
}
