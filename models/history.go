package models

import "time"

type History struct {
	ID         int       `db:"id"`
	UserID     int       `db:"user_id"`
	UserName   string    `db:"username"`
	CategoryID int       `db:"category_id"`
	ProjectID  int       `db:"project_id"`
	Flag       int       `db:"flag"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
