# lyticaa-api

Lyticaa API.

## Setup

### Golang

If you are new to Golang, please follow the setup instructions [here](https://golang.org/doc/install).

### Swagger

Please see the installation instructions [here](https://goswagger.io/install.html)

### Environment

Before running this project, please ensure that you have the following environment variables set:

```bash
APP_NAME=
SENTRY_DSN=
NEWRELIC_LICENSE_KEY=
PORT=
JWKS_URL=
JWT_AUD=
JWT_ISS=
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
AWS_REGION=
AWS_S3_UPLOAD_BUCKET=
```

If you are unsure as to what these values ought to be, then please check with a colleague.

### Linter

To run the linter:

```bash
make lint
```

### Tests

To run the tests and see test coverage:

```bash
make tests
```

### Install

To compile and install the binary:

```bash
make install
```

### Run the Service

```bash
make run-api-service
```

## Documentation 

### Swagger

To generate documentation for the API:

```bash
make generate-docs
```

The docs will then be available at: 

`./api/docs/index.html`

## Docker

A Docker stack is provided with this project. To boot the stack, simply run:

```bash
make run-stack
```

Please ensure that prior to running this, you add the above environment variables to the `build/.env` file. Docker Compose will use these when building the container.
