FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -installsuffix cgo -o app .


FROM scratch

COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /app
COPY --from=builder /app/app .

ENV GIN_MODE=release
ENV ENVIRONMENT=production
ARG exchange_app_id
ENV OPENEXCHANGERATES_APP_ID=$exchange_app_id

USER 1000:1000

EXPOSE 8080
CMD ["/app/app"]