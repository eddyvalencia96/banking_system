package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

var JWTKey = []byte("your_secret_key") // Cambia esto a una clave segura

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader != "" {
			tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return JWTKey, nil
			})

			if err != nil {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			if claims, ok := token.Claims.(*Claims); ok && token.Valid {
				// Almacenar las claims en el contexto de la solicitud para usarlas más adelante
				ctx := context.WithValue(r.Context(), "user", claims)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
		}
		http.Error(w, "Forbidden", http.StatusForbidden)
	})
}
