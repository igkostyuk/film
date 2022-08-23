package models

// Film base model
type Film struct {
	ID          int    `json:"film_id" db:"film_id"`
	Title       string `json:"title,omitempty" db:"title"`
	Description string `json:"description" db:"description"`
}
