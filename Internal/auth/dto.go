package auth

type ClientRegister struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email"     binding:"required,email"`
	Phone    string `json:"phone"     binding:"required"`
	Country  string `json:"country" binding:"omitempty iso3166_1_alpha2"`
	State    string `json:"state" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Password string `json:"password"  binding:"required,min=8,max=72"`
}
type ClientResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	State    string `json:"state"`
	Address  string `json:"address"`
}

type ShopRegister struct {
	BusinessName string   `json:"business_name" binding:"required"`
	Email        string   `json:"email"         binding:"required,email"`
	Phone        string   `json:"phone"         binding:"required"`
	Password     string   `json:"password"      binding:"required,min=8,max=72"`
	State        string   `json:"state" binding:"required"`
	Address      string   `json:"address"       binding:"required"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
}

type ShopResponse struct {
	ID           string  `json:"id"`
	BusinessName string  `json:"business_name"`
	Email        string  `json:"email"`
	Phone        string  `json:"phone"`
	State        string  `json:"state"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	IsVerified   bool    `json:"is_verified"`
}

type RiderRegister struct {
	FullName          string `json:"full_name"    binding:"required"`
	Email             string `json:"email"        binding:"required,email"`
	Phone             string `json:"phone"        binding:"required"`
	Password          string `json:"password"     binding:"required,min=8,max=72"`
	RegisteredCompany string `json:"registered_company" binding:"required"`
	PlateNumber       string `json:"plate_number"`
}
type RiderResponse struct {
	ID                string `json:"id"`
	FullName          string `json:"full_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	RegisteredCompany string `json:"registered_company"`
	IsVerified        bool   `json:"is_verified"`
}

type Login struct {
	Country  string `json:"country" binding:"omitempty, iso3166_1_alpha2"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password"   binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}
