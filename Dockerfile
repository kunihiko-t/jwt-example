FROM golang:1.8-alpine

RUN apk update && \
  apk --no-cache add git mercurial curl make gcc g++ bash

# RUN mkdir /
# WORKDIR /go
# ADD . /go
# RUN go get

# EXPOSE 8080
