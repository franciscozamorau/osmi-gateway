# Etapa de build
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copiar archivos de módulos
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Construir el binario del gateway
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway ./cmd/main.go

# Etapa de runtime
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/gateway .
COPY --from=builder /app/.env.production ./.env

EXPOSE 8080

CMD ["./gateway"]
