package api

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lyticaa/lyticaa-api/internal/types"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
)

func (a *Api) Start() {
	a.Logger.Info().Msgf("starting on %v....", ":"+os.Getenv("PORT"))
	a.Router.Use(a.forceSsl)

	a.Handlers()
	a.ErrorHandlers()

	a.Srv = &http.Server{
		Addr:         ":" + os.Getenv("PORT"),
		Handler:      a.Router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		if err := a.Srv.ListenAndServe(); err != nil {
			a.Logger.Info().Err(err)
		}
	}()
}

func (a *Api) Handlers() {
	a.Router.Use(a.commonMiddleware)

	a.healthCheck()
	a.usersUploadUrl(a.s3Client())
}

func (a *Api) ErrorHandlers() {
	a.Router.NotFoundHandler = handlers.LoggingHandler(
		os.Stdout,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
		),
	)

	a.Router.MethodNotAllowedHandler = handlers.LoggingHandler(
		os.Stdout,
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusMethodNotAllowed)
			},
		),
	)
}

func (a *Api) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Srv.Shutdown(ctx); err != nil {
		a.Logger.Fatal().Err(err)
	}

	a.Logger.Info().Msg("server exiting....")
}

func (a *Api) userIDFromClaim(r *http.Request) string {
	authHeaderParts := strings.Split(r.Header.Get("Authorization"), " ")
	tokenString := authHeaderParts[1]

	token, _ := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := a.getPemCert(token)
		if err != nil {
			a.Logger.Error().Err(err)
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})

	claims, _ := token.Claims.(*types.CustomClaims)
	parts := strings.Split(claims.UserID, "|")
	userID := parts[1]

	return userID
}
