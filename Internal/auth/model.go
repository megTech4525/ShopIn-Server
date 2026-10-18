package auth

import "time"

type Role string

const (
	RoleClient Role = "client"
	RoleShop   Role = "shop"
	RoleRider  Role = "rider"
)

type Account struct {
	ID        string
	Role      Role
	Emial     string
	Phone     string
	Password  string
	CreatedAt time.Time
}
