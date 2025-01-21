# Stage 1: Build the application
FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
# ARG RDS_DB_URL RDS_DB_USER RDS_DB_PASSWORD 
# ENV RDS_DB_URL=$RDS_DB_URL
# ENV RDS_DB_USER=$RDS_DB_USER
# ENV RDS_DB_PASSWORD=$RDS_DB_PASSWORD
RUN echo "Building application..."
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o server main.go

# Stage 2: Create the production-ready image
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY app.env .
COPY db/migration ./db/migration
COPY db/sqlc ./db/sqlc
COPY start.sh .
COPY wait-for.sh .
EXPOSE 8080
CMD [ "./server" ]
ENTRYPOINT [ "./start.sh" ]