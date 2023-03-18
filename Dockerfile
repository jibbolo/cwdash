# build stage
FROM golang:1.20.2 AS build-stage
ADD . /src/
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -o application

# run stage
FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /root/
COPY --from=build-stage /src/application application
ENV AWS_DEFAULT_REGION "us-east-2"
ENTRYPOINT ./application
