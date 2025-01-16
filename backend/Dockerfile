FROM golang:1.24rc1-alpine

WORKDIR /app

RUN go install github.com/mitranim/gow@latest

COPY backend/go.mod backend/go.sum ./
COPY .env .
RUN go mod download

COPY backend/google-service-account.json ./
ENV GOOGLE_APPLICATION_CREDENTIALS=/app/google-service-account.json

CMD ["gow", "run", "."]
