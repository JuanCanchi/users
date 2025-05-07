package domain

import "time"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // no devolver nunca el hash
	CreatedAt time.Time `json:"created_at"`
}
