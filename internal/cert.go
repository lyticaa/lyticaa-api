package api

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/lyticaa/lyticaa-api/internal/types"

	"github.com/dgrijalva/jwt-go"
)

func (a *Api) getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := a.Client.Get(os.Getenv("JWKS_URL"))

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = types.Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}
