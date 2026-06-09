FROM golang:1.26.3-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# -trimpath buat ngilangin path di binary
# -ldflags="-s -w" buat ngilangin debug info biar ukuran binary lebih kecil
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/kasvior-wallet ./cmd

FROM alpine:3.22

WORKDIR /app

# ca-certificates buat nanti kalo dah pake VPS, tzdata buat timezone datetime
RUN apk add --no-cache ca-certificates tzdata \
  && addgroup -S app \
  && adduser -S -G app app \
  && mkdir -p /app/public/img \
  && chown -R app:app /app

COPY --from=builder /out/kasvior-wallet /app/kasvior-wallet
COPY --from=builder /src/public /app/public

RUN chown -R app:app /app

USER app

EXPOSE 8080

ENTRYPOINT ["/app/kasvior-wallet"]
