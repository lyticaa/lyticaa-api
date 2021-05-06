package api

import (
	"encoding/json"
	"net/http"

	"github.com/lyticaa/lyticaa-api/internal/types"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/codegangsta/negroni"
)

func (a *Api) healthCheck() {
	a.Router.Handle("/api/health_check", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := types.Health{Status: "OK"}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			a.Logger.Error().Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jsonResponse)
		if err != nil {
			a.Logger.Error().Err(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))
}

func (a *Api) usersUploadUrl(client *s3.S3) {
	a.Router.Handle("/api/v1/users/upload_url", negroni.New(
		negroni.HandlerFunc(a.jwtMiddleware().HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := a.signedUploadUrl(a.userIDFromClaim(r), client)

			response := types.Url{Url: url}
			json, err := json.Marshal(response)
			if err != nil {
				a.Logger.Error().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(json)
			if err != nil {
				a.Logger.Error().Err(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}))),
	)
}
