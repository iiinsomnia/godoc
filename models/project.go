package models

import "time"

type Project struct {
	ID           int       `db:"id"`
	Name         string    `db:"name"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	Description  string    `db:"description"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
