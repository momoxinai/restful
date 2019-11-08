FROM library/golang:1.13 as builder


ARG SOURCE=/go/src/wwwin-github.cisco.com/DevNet/restful
ADD . ${SOURCE}
WORKDIR ${SOURCE}

ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -a -tags netgo -ldflags '-w' -o /go/bin/restful wwwin-github.cisco.com/DevNet/restful

FROM alpine:3.10

MAINTAINER DevNet Cloudy Team

LABEL Description="DevNet restful microservice image"

RUN apk update && \
    apk upgrade && \
    apk add \
        bash \
        ca-certificates \
    && rm -rf /var/cache/apk/*

COPY *.ini /restful/

ENV RUNMODE=stage

COPY --from=builder /go/bin/restful /restful/

WORKDIR /restful

ENTRYPOINT ["/restful/restful"]
