FROM golang:1.25.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o rapidlog ./cmd/app

FROM scratch
COPY --from=builder /app/rapidlog /rapidlog

CMD ["/rapidlog"]

