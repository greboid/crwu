package main

import (
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/greboid/rwtccus/containers"
	"github.com/greboid/rwtccus/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
)

var (
	WebPort     = flag.Int("web-port", 3000, "Port for webserver")
	Debug       = flag.Bool("debug", true, "Enable debug logging")
	InboundAuth = flag.String("in-token", "", "Token sent in the as auth bearer to validate request")
)

func main() {
	envflag.Parse()
	logger := createLogger(*Debug)
	log.Info().Msg("Starting rwtccus executor")
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(web.AuthMiddleware(*InboundAuth))
	r.Use(web.LoggerMiddleware(logger))
	r.Post("/", handle)
	log.Info().Str("URL", fmt.Sprintf("http://0.0.0.0:%d", *WebPort)).Msg("Starting webserver")
	ws := web.Web{}
	ws.Init(*WebPort, r)
	if err := ws.Run(); err != nil {
		log.Error().Err(err).Msg("error running server")
	}
	log.Info().Msg("Exiting")
}

func handle(w http.ResponseWriter, r *http.Request) {
	values := make([]string, 0)
	err := render.DecodeJSON(r.Body, &values)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "Unable to decode request"})
		log.Error().Err(err).Msg("Unable to decode request")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	go updateRequestedImages(values)
	render.Status(r, http.StatusOK)
}

func updateRequestedImages(images []string) {
	for index := range images {
		number, err := containers.UpdateMatchingContainers(images[index])
		if err != nil {
			log.Error().Err(err).Str("Image", images[index]).Msg("Unable to update containers")
			continue
		}
		if number > 0 {
			log.Info().Str("Image", images[index]).Int("Count", number).Msg("Updated containers")
		}
	}
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
