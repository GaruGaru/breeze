ARG GO_VERSION=1.14
ARG APP_NAME="breeze"

FROM golang:${GO_VERSION}-alpine AS builder

RUN  printf "nameserver 1.1.1.1\nnameserver 8.8.8.8\n" > /etc/resolv.conf && apk update && apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


ENTRYPOINT ["/app"]