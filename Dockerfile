FROM node:22-alpine AS web-builder

WORKDIR /app/web

COPY web/package.json web/package-lock.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=web-builder /app/web/dist ./web/dist
RUN \
  CGO_ENABLED=0 go build -o gate ./cmd/gate && \
  CGO_ENABLED=0 go build -o identity ./cmd/identity

FROM scratch AS gate
COPY --from=builder /app/gate /gate
EXPOSE 8080
CMD ["/gate"]

FROM scratch AS identity
COPY --from=builder /app/identity /identity
EXPOSE 50051
CMD ["/identity"]
