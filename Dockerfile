ARG GO_VERSION=1.14
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 && go build \
    -ldflags="-s -w" \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final

WORKDIR /
ENV PATH="/"

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app /breeze

ENTRYPOINT ["breeze"]