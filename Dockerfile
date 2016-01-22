FROM golang:1.5.3-alpine
MAINTAINER Raphael Randschau <nicolai86@me.com>

ADD . /go/src/github.com/nicolai86/unlocodes
RUN go install github.com/nicolai86/unlocodes


COPY . /src
WORKDIR /src

RUN apk update && apk upgrade && apk add curl
RUN curl -O http://www.unece.org/fileadmin/DAM/cefact/locode/loc152csv.zip

CMD []
ENTRYPOINT /go/bin/unlocodes

EXPOSE 3030
