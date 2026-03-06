// internal/middleware/rate_limit.go
package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter configuración para rate limiting
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // requests por segundo
	burst    int           // máximo de requests en ráfaga
	cleanup  time.Duration // tiempo para limpiar visitantes inactivos
}

type visitor struct {
	lastSeen time.Time
	tokens   float64
	mu       sync.Mutex
}

// NewRateLimiter crea una nueva instancia de rate limiter
func NewRateLimiter(rate, burst int, cleanup time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  cleanup,
	}
	go rl.cleanupVisitors()
	return rl
}

// Limit es el middleware de rate limiting
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener IP del cliente
		ip := r.RemoteAddr

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			v = &visitor{tokens: float64(rl.burst)}
			rl.visitors[ip] = v
		}
		rl.mu.Unlock()

		v.mu.Lock()
		defer v.mu.Unlock()

		// Calcular tokens disponibles
		now := time.Now()
		elapsed := now.Sub(v.lastSeen).Seconds()
		v.tokens += elapsed * float64(rl.rate)
		if v.tokens > float64(rl.burst) {
			v.tokens = float64(rl.burst)
		}
		v.lastSeen = now

		// Verificar si tiene tokens suficientes
		if v.tokens < 1 {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		v.tokens--
		next.ServeHTTP(w, r)
	})
}

// cleanupVisitors elimina visitantes inactivos
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(rl.cleanup)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
