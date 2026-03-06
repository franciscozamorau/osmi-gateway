package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
)

// JWTConfig contiene la configuración para el middleware JWT
type JWTConfig struct {
	SecretKey []byte
}

// NewJWTConfig crea una nueva configuración JWT
func NewJWTConfig(secret string) *JWTConfig {
	return &JWTConfig{
		SecretKey: []byte(secret),
	}
}

// Auth middleware para verificar tokens JWT
func (c *JWTConfig) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extraer token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Verificar formato Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validar token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validar método de firma
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return c.SecretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extraer claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Crear metadata con los claims para pasar al gRPC
		md := metadata.New(map[string]string{
			"user_id":    getStringClaim(claims, "user_id"),
			"user_email": getStringClaim(claims, "email"),
			"user_role":  getStringClaim(claims, "role"),
		})

		// Añadir metadata al contexto
		ctx := metadata.NewOutgoingContext(r.Context(), md)
		r = r.WithContext(ctx)

		// Continuar con la petición
		next.ServeHTTP(w, r)
	})
}

// OptionalAuth middleware que no rechaza si no hay token, solo añade metadata si existe
func (c *JWTConfig) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No hay token, continuar sin metadata
			next.ServeHTTP(w, r)
			return
		}

		// Verificar formato Bearer
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			// Formato inválido, continuar sin metadata
			next.ServeHTTP(w, r)
			return
		}

		tokenString := parts[1]

		// Validar token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return c.SecretKey, nil
		})

		if err != nil || !token.Valid {
			// Token inválido, continuar sin metadata
			next.ServeHTTP(w, r)
			return
		}

		// Extraer claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			next.ServeHTTP(w, r)
			return
		}

		// Crear metadata con los claims
		md := metadata.New(map[string]string{
			"user_id":    getStringClaim(claims, "user_id"),
			"user_email": getStringClaim(claims, "email"),
			"user_role":  getStringClaim(claims, "role"),
		})

		ctx := metadata.NewOutgoingContext(r.Context(), md)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// Helper para extraer claims string de forma segura
func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// ExtractUserIDFromContext extrae el user_id del contexto (para usar en handlers)
func ExtractUserIDFromContext(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	values := md.Get("user_id")
	if len(values) == 0 {
		return "", false
	}

	return values[0], true
}

// RequireRole middleware para verificar roles específicos
func (c *JWTConfig) RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer role del contexto (asumimos que Auth middleware ya corrió)
			md, ok := metadata.FromIncomingContext(r.Context())
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			roles := md.Get("user_role")
			if len(roles) == 0 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			userRole := roles[0]

			// Verificar si el rol del usuario está en los roles permitidos
			for _, role := range allowedRoles {
				if role == userRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
		})
	}
}
