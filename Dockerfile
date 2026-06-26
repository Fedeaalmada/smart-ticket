# ── Etapa 1: compilar el backend Go ──────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar dependencias primero (mejor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o smartticket ./main.go

# ── Etapa 2: imagen final liviana ────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Copiar solo el binario compilado
COPY --from=builder /app/smartticket .

EXPOSE 8080

CMD ["./smartticket"]
