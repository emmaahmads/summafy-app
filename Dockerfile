# Stage 1: Build the application
FROM golang:1.23.4-alpine3.20 AS builder
WORKDIR /app
RUN echo "Building application..."
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o server main.go

# Stage 2: Create the production-ready image
FROM alpine:3.20
WORKDIR /app
RUN adduser -D -g '' appuser && chown appuser /app
USER appuser
COPY --from=builder /app/server .
COPY app.env .
COPY db/migration ./db/migration
COPY db/sqlc ./db/sqlc
COPY start.sh .
COPY wait-for.sh .
EXPOSE 8080
CMD [ "./server" ]
ENTRYPOINT [ "./start.sh" ]