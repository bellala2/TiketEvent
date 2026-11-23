package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bellala2/TiketEvent/internal/response"
	"github.com/bellala2/TiketEvent/internal/service"
)

// AuthUser adalah identitas user yang sudah lolos pemeriksaan token.
type AuthUser struct {
	ID   string
	Role string
}

// ctxKey bertipe khusus supaya tidak bentrok dengan key milik package lain.
type ctxKey struct{}

// TokenParser cukup berupa interface kecil agar mudah ditiru saat test.
type TokenParser interface {
	Parse(token string) (*service.Claims, error)
}

// UserFromContext dipakai handler untuk mengetahui siapa yang sedang request.
func UserFromContext(ctx context.Context) (AuthUser, bool) {
	u, ok := ctx.Value(ctxKey{}).(AuthUser)
	return u, ok
}

var errUnauthorized = service.NewError(service.CodeUnauthorized, "Token tidak valid atau sudah kedaluwarsa")

// Authenticate memeriksa header "Authorization: Bearer <token>".
// Token tidak ada, rusak, atau kedaluwarsa -> 401.
func Authenticate(tokens TokenParser) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			scheme, tokenString, found := strings.Cut(header, " ")
			if !found || !strings.EqualFold(scheme, "Bearer") || strings.TrimSpace(tokenString) == "" {
				response.WriteError(w, errUnauthorized)
				return
			}

			claims, err := tokens.Parse(strings.TrimSpace(tokenString))
			if err != nil {
				response.WriteError(w, errUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, AuthUser{ID: claims.Subject, Role: claims.Role})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole membatasi endpoint untuk role tertentu.
// Harus dipasang SETELAH Authenticate. Role tidak cocok -> 403.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				response.WriteError(w, errUnauthorized)
				return
			}
			for _, role := range roles {
				if user.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			response.WriteError(w, service.NewError(service.CodeForbidden, "Anda tidak memiliki akses ke resource ini"))
		})
	}
}
