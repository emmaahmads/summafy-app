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

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server .
COPY app.env .
COPY templates ./templates
COPY db/migration ./db/migration
COPY db/sqlc ./db/sqlc
COPY start.sh .
COPY wait-for.sh .
EXPOSE 8080
CMD [ "./server" ]
ENTRYPOINT [ "./start.sh" ]

# # Stage 2: Create the production-ready image
# FROM alpine:3.20
# WORKDIR /app
# COPY --from=builder /app/summafy .

# COPY main.go .
# COPY app.env .
# COPY Makefile .
# RUN go get -u github.com/emmaahmads/summafy/util
# RUN go get -u github.com/lib/pq
# RUN go get -u github.com/sirupsen/logrus
# RUN go get -u github.com/gin-gonic/gin
# # RUN go get -u github.com/emmaahmads/summafy/db/sqlc
# # RUN go get -u github.com/emmaahmads/summafy/api
# RUN go get -u github.com/aws/aws-sdk-go-v2/aws
# RUN go get -u github.com/aws/aws-sdk-go-v2/aws/session
# RUN go get -u github.com/aws/aws-sdk-go-v2/service/s3
# RUN go get -u github.com/sashabaranov/go-openai
# RUN go get -u github.com/gin-contrib/sessions
