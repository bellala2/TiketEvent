package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/bellala2/TiketEvent/internal/clock"
)

// Claims adalah isi (payload) token: sub (user ID), role, iat, exp.
type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
	clock  clock.Clock
}

func NewTokenManager(secret string, ttl time.Duration, clk clock.Clock) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl, clock: clk}
}

func (m *TokenManager) TTL() time.Duration { return m.ttl }

// Generate membuat token JWT (HS256) untuk user.
func (m *TokenManager) Generate(userID, role string) (string, error) {
	now := m.clock.Now()
	claims := Claims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Parse memeriksa tanda tangan dan masa berlaku token.
// Hanya algoritma HS256 yang diterima (mencegah serangan "alg: none").
func (m *TokenManager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(t *jwt.Token) (any, error) { return m.secret, nil },
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithTimeFunc(m.clock.Now),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid || claims.Subject == "" {
		return nil, fmt.Errorf("token tidak valid")
	}
	return claims, nil
}
