FROM golang:1.5.3-alpine
MAINTAINER Raphael Randschau <nicolai86@me.com>

ADD . /go/src/github.com/nicolai86/unlocodes
RUN go install github.com/nicolai86/unlocodes

EXPOSE 3030

COPY . /src
WORKDIR /src

ENTRYPOINT /go/bin/unlocodes
