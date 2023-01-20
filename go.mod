module github.com/lyticaa/lyticaa-api

replace github.com/lyticaa/lyticaa-api => ../lyticaa/lyticaa-api

go 1.13

require (
	bou.ke/monkey v1.0.2
	github.com/auth0/go-jwt-middleware v0.0.0-20190805220309-36081240882b
	github.com/aws/aws-sdk-go v1.44.183
	github.com/codegangsta/negroni v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/sentry-go v0.7.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/johannesboyne/gofakes3 v0.0.0-20191029185751-e238f04965fe
	github.com/newrelic/go-agent v2.16.3+incompatible
	github.com/rs/zerolog v1.20.0
)
