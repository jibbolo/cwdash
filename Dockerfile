# build stage
FROM golang:1.19 AS build-stage
ARG VERSION
ADD . /src/
WORKDIR /src/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.VERSION=$VERSION -X main.BUILDDATE=`date -u +%Y%m%d.%H%M%S`" -o application

# run stage
FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /root/

COPY --from=build-stage /src/application application
COPY templates templates

ENV AWS_DEFAULT_REGION "us-west-2"

ENTRYPOINT ./application

EXPOSE 5001
