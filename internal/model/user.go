package model

import (
	"database/sql"
	"time"
)

type UserFilter struct {
	UserID int `json:"Id"`
}

type User struct {
	ID              int    `db:"id"`
	Name            string `db:"name"`
	Email           string `db:"email"`
	Role            int    `db:"role"`
	Password        string `db:"password"`
	ConfirmPassword string `db:"password_confirm"`

	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
