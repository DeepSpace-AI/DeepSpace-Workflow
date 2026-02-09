FROM golang:1.25.6-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY services/worker/go.mod services/worker/go.sum ./services/worker/
WORKDIR /src/services/worker
RUN go mod download

COPY services/worker ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/worker ./cmd/worker

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/worker /app/worker
COPY templates /app/templates

CMD ["/app/worker"]
