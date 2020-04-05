ARG GO_VERSION=1.14
ARG APP_NAME="breeze"

FROM golang:${GO_VERSION} AS builder

RUN apt-get update && apt-get install -y ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w" \
    -installsuffix 'static' \
    -o /app .

FROM scratch AS final
ENV PATH="/:${PATH}"
COPY --from=builder /app /breeze
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["breeze"]