#
# Build
#
FROM golang:1.13 AS build

WORKDIR /app

COPY . .

ENV GOBIN=/go/bin
ENV GOPATH=/go
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go mod verify
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o $GOBIN/gapid ./cmd/apid

#
# Main
#
FROM alpine

# Copy the compiled binaires over.
COPY --from=build /go/bin/apid /usr/bin/

EXPOSE 3000
