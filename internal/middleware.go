package api

import (
	"errors"
	"net/http"
	"os"

	jwtm "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
)

func (a *Api) jwtMiddleware() *jwtm.JWTMiddleware {
	return jwtm.New(jwtm.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			check := token.Claims.(jwt.MapClaims).VerifyAudience(os.Getenv("JWT_AUD"), false)
			if !check {
				return token, errors.New("invalid audience")
			}

			check = token.Claims.(jwt.MapClaims).VerifyIssuer(os.Getenv("JWT_ISS"), false)
			if !check {
				return token, errors.New("invalid issuer")
			}

			cert, err := a.getPemCert(token)
			if err != nil {
				a.Logger.Error().Err(err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

func (a *Api) commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (a *Api) forceSsl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENV") != "development" {
			if r.Header.Get("x-forwarded-proto") != "https" {
				sslUrl := "https://" + r.Host + r.RequestURI
				http.Redirect(w, r, sslUrl, http.StatusTemporaryRedirect)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
