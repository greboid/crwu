package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/csmith/envflag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/greboid/rwtccus/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	WebPort      = flag.Int("web-port", 3000, "Port for webserver")
	Debug        = flag.Bool("debug", true, "Enable debug logging")
	InboundToken = flag.String("in-token", "", "Token sent in the auth bearer to validate requests from the registry")
	OutboundAuth = flag.String("out-token", "", "Token sent in the auth bearer to validate requests to the executor")
	ExecutorURL  = flag.String("executor", "", "URL to send executor requests to")
)

func main() {
	envflag.Parse()
	logger := createLogger(*Debug)
	log.Info().Msg("Starting rwtccus listener")
	if ExecutorURL == nil || *ExecutorURL == "" {
		log.Fatal().Msg("Executor URL is required")
	}
	if _, err := url.Parse(*ExecutorURL); err != nil {
		log.Fatal().Err(err).Msg("Unable to parse executor")
	}
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(web.AuthMiddleware(*InboundToken))
	r.Use(web.LoggerMiddleware(logger))
	r.Post("/distribution", handleDistributionWebhook)
	log.Info().Str("URL", fmt.Sprintf("http://0.0.0.0:%d", *WebPort)).Msg("Starting webserver")
	ws := web.Web{}
	ws.Init(*WebPort, r)
	if err := ws.Run(); err != nil {
		log.Error().Err(err).Msg("error running server")
	}
	log.Info().Msg("Exiting")
}

func handleDistributionWebhook(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()
	if err != nil {
		log.Error().Err(err).Msg("Error reading body")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"message": "OK"})
		return
	}
	webhook := &Webhook{}
	err = render.DecodeJSON(bytes.NewReader(bodyBytes), webhook)
	if err != nil {
		//Whilst we can't process this, send an OK response to the registry
		//Webhooks backup on the registry if they're not accepted causing issues
		log.Error().Err(err).Msg("Unable to decode webhook")
		render.Status(r, http.StatusOK)
		return
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "OK"})
	go sendRequestToExecutor(webhook)
}

func sendRequestToExecutor(webhook *Webhook) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	var images []string
	for index := range webhook.Events {
		switch webhook.Events[index].Action {
		case "pull":
			//Discard
			break
		case "push":
			log.Debug().Interface("Webhook event", webhook.Events[index]).Msg("Push received")
			images = append(images, fmt.Sprintf("%s/%s:%s", webhook.Events[index].Request.Host, webhook.Events[index].Target.Repository, webhook.Events[index].Target.Tag))
			break
		default:
			log.Debug().Interface("Webhook event", webhook.Events[index]).Msg("Unknown webhook action")
		}
	}
	if len(images) == 0 {
		return
	}
	jsonBytes, err := json.Marshal(images)
	if err != nil {
		log.Error().Err(err).Msg("Unable to create request to executor")
	}
	req, err := http.NewRequest(http.MethodPost, *ExecutorURL, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Error().Err(err).Msg("Unable to create request to executor")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *OutboundAuth))
	res, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Unable to send request to executor")
	}
	if res.StatusCode < 200 && res.StatusCode <= 300 {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Err(err).Msg("Unable to read error response")
		}
		log.Error().Int("Status", res.StatusCode).Str("Response", string(bodyBytes)).Msg("Error response when sending request to executor")
		return
	}
	log.Info().Strs("Images", images).Msg("Update request sent")
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
