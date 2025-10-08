package middleware

import (
	"context"
	"fmt"
	"gonotes/internal/api"
	"gonotes/internal/security"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{secret: []byte(secret)}
}

func (m *AuthMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("missing or invalid token")))
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &security.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return m.secret, nil
		})

		if err != nil || !token.Valid {
			render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("invalid token")))
			return
		}

		userID, err := strconv.Atoi(claims.UserID)
		if err != nil {
			render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("invalid user id")))
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
