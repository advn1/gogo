package user

import "github.com/google/uuid"

// User scheme for in-database model
type User struct {
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password_hash"`
	Id            uuid.UUID `json:"id"`
}
