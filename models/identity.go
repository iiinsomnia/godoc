package models

import "time"

type Identity struct {
	ID            int       `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Email         string    `db:"email" json:"email"`
	Password      string    `db:"password" json:"password"`
	Salt          string    `db:"salt" json:"salt"`
	Role          int       `db:"role" json:"role"`
	LastLoginIP   string    `db:"last_login_ip" json:"last_login_ip"`
	LastLoginTime time.Time `db:"last_login_time" json:"last_login_time"`
}
