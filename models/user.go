package models

import "time"

type User struct {
	ID            int       `db:"id"`
	Name          string    `db:"name"`
	Email         string    `db:"email"`
	Role          int       `db:"role"`
	LastLoginIP   string    `db:"last_login_ip"`
	LastLoginTime time.Time `db:"last_login_time"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
