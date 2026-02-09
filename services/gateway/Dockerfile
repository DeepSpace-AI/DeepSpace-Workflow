FROM golang:1.25.6-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/
WORKDIR /src/services/gateway
RUN go mod download

COPY services/gateway ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/gateway ./cmd/gateway
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/admin-init ./cmd/admin-init

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /out/gateway /app/gateway
COPY --from=builder /out/admin-init /app/admin-init
COPY templates /app/templates

RUN mkdir -p /data/kb

EXPOSE 8080

CMD ["/app/gateway"]
