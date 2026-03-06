// internal/routes/router.go
package routes

import (
	"net/http"
	"os"

	"github.com/franciscozamorau/osmi-gateway/internal/handlers"
	"github.com/franciscozamorau/osmi-gateway/internal/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func SetupRouter(gwmux *runtime.ServeMux) http.Handler {
	// Configuración JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "tu-secreto-jwt-por-defecto-cambiar-en-produccion"
	}
	authMiddleware := middleware.NewJWTConfig(jwtSecret)

	// Rutas públicas (no requieren autenticación)
	gwmux.HandlePath("GET", "/health", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		handlers.HealthHandler(w, r)
	})

	// NOTA: Las rutas gRPC (/customers, /events, /tickets, etc.)
	// ya están registradas automáticamente por RegisterOsmiServiceHandlerFromEndpoint
	// en main.go. NO necesitas duplicarlas aquí.

	// Ruta protegida de ejemplo (opcional)
	protectedHandler := authMiddleware.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Acceso autorizado a ruta protegida"}`))
	}))

	gwmux.HandlePath("GET", "/api/protected", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		protectedHandler.ServeHTTP(w, r)
	})

	// Aplicar middleware a TODAS las rutas
	handler := middleware.Logging(gwmux)
	handler = middleware.CORS(handler)

	// Rate limiting (opcional, descomentar si lo necesitas)
	// rateLimiter := middleware.NewRateLimiter(10, 20, 3*time.Minute)
	// handler = rateLimiter.Limit(handler)

	return handler
}
