package web

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
)

func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			wrapper := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				logger.Info().
					Timestamp().
					Str("REMOTE_ADDR", r.RemoteAddr).
					Str("url", r.URL.Path).
					Str("Proto", r.Proto).
					Str("Method", r.Method).
					Str("User-Agent", r.Header.Get("User-Agent")).
					Int("status", wrapper.Status()).
					Msg("incoming_request")
			}()
			next.ServeHTTP(wrapper, r)
		})
	}
}

func AuthMiddleware(authToken string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				render.Status(r, http.StatusUnauthorized)
			}
			splitAuth := strings.Split(auth, "Bearer")
			if len(splitAuth) != 2 {
				render.Status(r, http.StatusUnauthorized)
			}
			suppliedToken := strings.TrimSpace(splitAuth[1])
			if len(suppliedToken) < 1 {
				render.Status(r, http.StatusUnauthorized)
			}
			if authToken != suppliedToken {
				render.Status(r, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}
