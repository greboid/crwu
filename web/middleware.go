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
				logger.Trace().
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
			suppliedToken := getAuthHeader(r.Header.Get("Authorization"))
			if authToken != suppliedToken {
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, map[string]string{"error": "unauthorised"})
				return
			}
			next.ServeHTTP(w, r)
			return
		})
	}
}

func getAuthHeader(authHeader string) string {
	if authHeader == "" {
		return ""
	}
	splitAuth := strings.Split(authHeader, "Bearer")
	if len(splitAuth) != 2 {
		return ""
	}
	token := strings.TrimSpace(splitAuth[1])
	if len(token) < 1 {
		return ""
	}
	return token
}
