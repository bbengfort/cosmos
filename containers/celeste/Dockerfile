# Dynamic Builds
ARG BUILDER_IMAGE=golang:1.17-buster
ARG FINAL_IMAGE=debian:buster-slim

# Build stage
FROM ${BUILDER_IMAGE} as builder

# Ensure ca-certificates are up to date
RUN update-ca-certificates

WORKDIR $GOPATH/src/github.com/bbengfort/cosmos

# Use modules for dependencies
COPY go.mod .
COPY go.sum .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

# Copy only what is needed (ignore web and fixtures)
COPY cmd ./cmd
COPY pkg ./pkg
COPY proto ./proto

# Run tests to avoid building a broken container
RUN go test -timeout 180s ./...

# Build binaries
RUN go build -v -o /go/bin/celeste ./cmd/celeste
RUN go build -v -o /go/bin/cosmos ./cmd/cosmos

# Final Stage
FROM ${FINAL_IMAGE} AS final

LABEL maintainer="Benjamin Bengfort <support@bengfort.com>"
LABEL description="Celeste: The Cosmos Game Server"

# Ensure ca-certificates are up to date
RUN set -x && apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/bin/celeste /usr/local/bin/celeste
COPY --from=builder /go/bin/cosmos /usr/local/bin/cosmos

CMD [ "/usr/local/bin/celeste", "serve" ]