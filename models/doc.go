package models

import "time"

type Doc struct {
	ID           int       `db:"id"`
	Title        string    `db:"title"`
	CategoryID   int       `db:"category_id"`
	CategoryName string    `db:"category_name"`
	ProjectID    int       `db:"project_id"`
	ProjectName  string    `db:"project_name"`
	Label        string    `db:"label"`
	Markdown     string    `db:"markdown"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
