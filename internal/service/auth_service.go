package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"

	"github.com/bellala2/TiketEvent/internal/model"
	"github.com/bellala2/TiketEvent/internal/repository"
)

// UserRepository didefinisikan di sisi pemakai (service), bukan di repository.
// Dengan begitu saat unit test kita cukup membuat repository tiruan.
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}

// Pesan login SENGAJA sama untuk email tak terdaftar maupun password salah,
// supaya orang tidak bisa menebak email mana yang terdaftar.
var errInvalidCredentials = NewError(CodeUnauthorized, "Email atau password salah")

type AuthService struct {
	users  UserRepository
	tokens *TokenManager
}

func NewAuthService(users UserRepository, tokens *TokenManager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.RegisteredUser, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = normalizeEmail(req.Email)

	if details := validateRegister(req); len(details) > 0 {
		return nil, NewValidationError("Validasi gagal", details...)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         model.RoleCustomer, // register SELALU customer
	}
	if err := s.users.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, NewError(CodeConflict, "Email sudah terdaftar")
		}
		return nil, err
	}

	return &model.RegisteredUser{
		UserResponse: toUserResponse(user),
		CreatedAt:    user.CreatedAt.UTC().Format(time.RFC3339),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	email := normalizeEmail(req.Email)

	var details []model.FieldError
	if email == "" {
		details = append(details, model.FieldError{Field: "email", Message: "wajib diisi"})
	}
	if req.Password == "" {
		details = append(details, model.FieldError{Field: "password", Message: "wajib diisi"})
	}
	if len(details) > 0 {
		return nil, NewValidationError("Validasi gagal", details...)
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, errInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errInvalidCredentials
	}

	token, err := s.tokens.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.tokens.TTL().Seconds()),
		User:        toUserResponse(user),
	}, nil
}

// GetMe mengambil data user dari ID di dalam token.
func (s *AuthService) GetMe(ctx context.Context, userID string) (*model.UserResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Token valid tapi usernya sudah tidak ada: anggap tidak terautentikasi.
			return nil, NewError(CodeUnauthorized, "Token tidak valid atau sudah kedaluwarsa")
		}
		return nil, err
	}
	resp := toUserResponse(user)
	return &resp, nil
}

// EnsureAdmin membuat akun admin pertama kalau belum ada (seed saat aplikasi start).
// Aman dijalankan berulang kali dan aman kalau ada lebih dari satu instance.
func (s *AuthService) EnsureAdmin(ctx context.Context, email, password string) error {
	email = normalizeEmail(email)

	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return nil // sudah ada
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password admin: %w", err)
	}

	err = s.users.Create(ctx, &model.User{
		Name:         "Administrator",
		Email:        email,
		PasswordHash: string(hash),
		Role:         model.RoleAdmin,
	})
	if errors.Is(err, repository.ErrDuplicateEmail) {
		return nil // instance lain sudah membuatnya lebih dulu
	}
	return err
}

func toUserResponse(u *model.User) model.UserResponse {
	return model.UserResponse{ID: u.ID, Name: u.Name, Email: u.Email, Role: u.Role}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateRegister(req model.RegisterRequest) []model.FieldError {
	var details []model.FieldError

	if n := utf8.RuneCountInString(req.Name); n < 3 || n > 100 {
		details = append(details, model.FieldError{Field: "name", Message: "harus 3-100 karakter"})
	}

	if !isValidEmail(req.Email) {
		details = append(details, model.FieldError{Field: "email", Message: "format email tidak valid"})
	}

	// bcrypt hanya membaca 72 byte pertama, jadi batas atas dihitung dalam byte.
	if len(req.Password) < 8 || len(req.Password) > 72 {
		details = append(details, model.FieldError{Field: "password", Message: "harus 8-72 karakter"})
	} else if !hasLetterAndDigit(req.Password) {
		details = append(details, model.FieldError{Field: "password", Message: "minimal satu huruf dan satu angka"})
	}

	return details
}

func isValidEmail(email string) bool {
	if email == "" || len(email) > 150 {
		return false
	}
	addr, err := mail.ParseAddress(email)
	// Address harus sama persis, supaya input seperti "Budi <b@x.com>" ditolak.
	return err == nil && addr.Address == email
}

func hasLetterAndDigit(s string) bool {
	var letter, digit bool
	for _, r := range s {
		switch {
		case unicode.IsLetter(r):
			letter = true
		case unicode.IsDigit(r):
			digit = true
		}
	}
	return letter && digit
}
