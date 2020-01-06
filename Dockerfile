ARG GO_VERSION=1.13.5
FROM golang:${GO_VERSION}-alpine AS builder

ARG ARCH=amd64
ARG OS=linux

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN export GOARCH=$ARCH GOOS=$OS CGO_ENABLED=0 && go build \
    -ldflags="-s -w" \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]