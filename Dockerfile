# Stage 1: Build the application
FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY templates templates
COPY db/migration db/migration
EXPOSE 8080
CMD [ "/app/main" ]

# Stage 2: Create the production-ready image
# FROM alpine:3.20
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY templates templates
# COPY db/migration db/migration
# COPY db/sqlc db/sqlc
# COPY api/api.go api/
# COPY api/aws_s3.go api/
# COPY api/dashboard.go api/
# COPY api/ai.go api/
# COPY main.go .
# COPY app.env .
# COPY Makefile .
# RUN go get -u github.com/emmaahmads/summafy/util
# RUN go get -u github.com/lib/pq
# RUN go get -u github.com/sirupsen/logrus
# RUN go get -u github.com/gin-gonic/gin
# RUN go get -u github.com/emmaahmads/summafy/db/sqlc
# RUN go get -u github.com/emmaahmads/summafy/api
# RUN go get -u github.com/aws/aws-sdk-go/aws
# RUN go get -u github.com/aws/aws-sdk-go/aws/session
# RUN go get -u github.com/aws/aws-sdk-go/service/s3
# RUN go get -u github.com/sashabaranov/go-openai
# RUN go get -u github.com/gin-contrib/sessions
# EXPOSE 8080
# CMD ["./main"]