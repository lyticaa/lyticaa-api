package api

import (
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Api struct {
	Logger   zerolog.Logger
	NewRelic newrelic.Application
	Srv      *http.Server
	Router   *mux.Router
	Client   *http.Client
}

func NewApi() *Api {
	sentryOpts := sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	}
	err := sentry.Init(sentryOpts)
	if err != nil {
		panic(err)
	}

	config := newrelic.NewConfig(
		os.Getenv("APP_NAME"),
		os.Getenv("NEWRELIC_LICENSE_KEY"),
	)
	nr, _ := newrelic.NewApplication(config)

	return &Api{
		Logger:   log.With().Str("module", os.Getenv("APP_NAME")).Logger(),
		NewRelic: nr,
		Router:   mux.NewRouter(),
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
