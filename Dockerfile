FROM golang:alpine

RUN apk apk update && apk upgrade && apk add --no-cache git ruby-rake

ADD . /muffet
WORKDIR /muffet

RUN go get -d
RUN rake build

RUN apk del git ruby-rake

ENTRYPOINT ["/muffet/muffet"]
