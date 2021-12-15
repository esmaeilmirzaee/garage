package user

import (
	"github.com/lib/pq"
	"time"
)

type User struct {
	ID        string         `db:"user_id" json:"user_id"`
	Name      string         `db:"name" json:"name"`
	Email     string         `db:"email" json:"email"`
	Password  []byte         `db:"password" json:"-"`
	Roles     pq.StringArray `db:"roles" json:"roles"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required"`
	Roles           []string `json:"roles" validate:"required"`
	Password        string   `json:"password" validate:"required"`
	ConfirmPassword string   `json:"confirm_password" validate:"eqfield=Password"`
}
