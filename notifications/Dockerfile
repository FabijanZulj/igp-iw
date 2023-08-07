FROM golang:1.18-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o notifications ./cmd/main.go

FROM alpine:latest

COPY --from=builder ["/build/notifications", "/"]
COPY --from=builder ["/build/templates", "/templates"]
COPY --from=builder ["/build/app.env", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/notifications"]
