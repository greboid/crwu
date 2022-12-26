package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
)

var (
	WebPort = flag.Int("web-port", 3000, "Port for webserver")
	Debug   = flag.Bool("debug", true, "Enable debug logging")
)

func main() {
	envflag.Parse()
	logger := createLogger(*Debug)
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(loggerMiddleware(logger))
	r.Get("/", handleIndex)
	r.Get("/hook", handleWebhook)
	server := http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *WebPort),
		Handler: r,
	}
	log.Info().Str("URL", fmt.Sprintf("http://0.0.0.0:%d", *WebPort)).Msg("Starting webserver")
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, os.Kill)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Error with server")
		}
		log.Info().Msg("Server shutdown")
	}()
	<-signals
	_ = server.Shutdown(context.Background())
	log.Info().Msg("Exiting")
}

func handleIndex(_ http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusOK)
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	var webhook *Webhook
	err := render.DecodeJSON(r.Body, webhook)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Bad input"})
		return
	}
	render.Status(r, http.StatusOK)
}

func createLogger(debug bool) *zerolog.Logger {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	return &logger
}

func loggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
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
