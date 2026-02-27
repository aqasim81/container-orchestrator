package api

import (
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/github-builder/container-orchestrator/internal/api/handlers"
	"github.com/github-builder/container-orchestrator/internal/store"
)

// RouterConfig holds dependencies for constructing the API router.
type RouterConfig struct {
	Store        store.Store
	Logger       zerolog.Logger
	DashboardURL string
	APIKey       string `json:"-"` //nolint:gosec // Not a hardcoded credential, populated from env.
}

// NewRouter creates a Chi router with middleware and all API routes mounted.
func NewRouter(cfg *RouterConfig) http.Handler {
	r := chi.NewRouter()

	// Middleware chain.
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(zerologMiddleware(&cfg.Logger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.DashboardURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-API-Key", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Health endpoint — no auth required.
	r.Get("/healthz", handlers.Health())

	// API v1 routes — auth required.
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(apiKeyAuth(cfg.APIKey))
		// Future route groups will be mounted here.
	})

	return r
}

// zerologMiddleware logs each HTTP request with zerolog.
func zerologMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info().
					Str("request_id", middleware.GetReqID(r.Context())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", ww.Status()).
					Int("bytes", ww.BytesWritten()).
					Dur("duration", time.Since(start)).
					Msg("request completed")
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

// apiKeyAuth validates the X-API-Key header.
func apiKeyAuth(apiKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" || subtle.ConstantTimeCompare([]byte(key), []byte(apiKey)) != 1 {
				Error(w, http.StatusUnauthorized, "invalid or missing API key", "UNAUTHORIZED")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
