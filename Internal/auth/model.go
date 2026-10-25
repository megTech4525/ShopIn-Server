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
type ClientProfile struct {
	FullName string
	State    string
	Address  string
}

type ShopProfile struct {
	BusinessName string
	State        string
	Address      string
	Latitude     *float64
	Longitude    *float64
}

type RiderProfile struct {
	FullName          string
	RegisteredCompany string
	PlateNumber       *string
}
