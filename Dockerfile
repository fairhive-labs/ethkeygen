# ---- build stage ----
# golang:alpine bundles git and, with CGO disabled, needs no C toolchain.
FROM golang:1.24-alpine AS builder
WORKDIR /src

# Download dependencies in a separate, cacheable layer.
COPY go.mod go.sum ./
RUN go mod download

# Build a static, stripped binary (no libc dependency -> runs on bare alpine).
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./api

# ---- runtime stage ----
FROM alpine:3.21
RUN apk add --no-cache ca-certificates wget \
    && adduser -D -u 10001 appuser
COPY --from=builder /out/api /app/bin/api
USER appuser
ENV GIN_MODE=release
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget -qO- http://localhost:8080/health || exit 1
CMD ["/app/bin/api"]
