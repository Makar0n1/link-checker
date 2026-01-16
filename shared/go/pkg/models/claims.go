package models

import "github.com/golang-jwt/jwt/v5"

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Claims represents JWT token claims
type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   Role   `json:"role"`
	jwt.RegisteredClaims
}
