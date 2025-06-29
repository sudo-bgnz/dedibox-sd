# syntax=docker/dockerfile:1

# --- Builder stage ---
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go mod download
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /app/scw_sd main.go

# --- Runtime stage ---
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
COPY --from=build /app/scw_sd /usr/local/bin/scw_sd
EXPOSE 8000
ENTRYPOINT ["/usr/local/bin/scw_sd"]
