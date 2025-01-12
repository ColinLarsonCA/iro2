FROM golang:1.24rc1-alpine

WORKDIR /app

RUN go install github.com/mitranim/gow@latest

COPY backend/go.mod backend/go.sum ./
COPY .env .
RUN go mod download

CMD ["gow", "run", "."]
