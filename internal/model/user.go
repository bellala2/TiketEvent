package model

import "time"

const (
	RoleAdmin    = "admin"
	RoleCustomer = "customer"
)

// User adalah representasi baris tabel users.
// PasswordHash diberi tag json:"-" supaya TIDAK PERNAH ikut terkirim ke client.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse adalah data user yang aman dikirim ke client.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// RegisteredUser dipakai di response register (ada created_at).
// UserResponse yang di-embed otomatis "diratakan" di JSON.
type RegisteredUser struct {
	UserResponse
	CreatedAt string `json:"created_at"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int          `json:"expires_in"`
	User        UserResponse `json:"user"`
}
