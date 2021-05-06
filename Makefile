GOBIN?=${GOPATH}/bin

all: lint install

lint-pre:
	@test -z $(gofmt -l .)
	@go mod verify

lint: lint-pre
	@golangci-lint run

lint-verbose: lint-pre
	@golangci-lint run -v

install: go.sum
	GO111MODULE=on go install -v ./cmd/apid

clean:
	rm -f ${GOBIN}/{apid}

tests:
	@go test -mod=readonly -v -coverprofile .testCoverage.txt ./...

setup-yarn:
	yarn install

run-api-service:
	@apid

run-stack:
	@docker-compose -f ./build/docker-compose.yml up --force-recreate --remove-orphans

generate-docs: setup-yarn
	./node_modules/.bin/redoc-cli bundle ./api/docs/openapi.yml -o ./api/docs/index.html
