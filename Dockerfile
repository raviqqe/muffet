FROM golang:alpine

RUN apk apk update && apk upgrade && apk add --no-cache git

ADD . /muffet
WORKDIR /muffet

RUN go get -d
RUN go build

ENTRYPOINT ["/muffet/muffet"]
