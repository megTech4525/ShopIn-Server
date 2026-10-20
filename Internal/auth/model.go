package auth

import "time"

type Role string

const (
	RoleClient Role = "client"
	RoleShop   Role = "shop"
	RoleRider  Role = "rider"
)

type Account struct {
	ID           string
	Role         Role
	Email        string
	Phone        string
	PasswordHash string
	CreatedAt    time.Time
}
