FROM golang:1.18-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserver ./cmd/main.go

FROM scratch

COPY --from=builder ["/build/apiserver", "/build/app.env", "/build/migrations", "/"]
COPY --from=builder ["/build/app.env", "/"]
COPY --from=builder ["/build/migrations", "/migrations"]

ENTRYPOINT ["/apiserver"]

