FROM golang:alpine

RUN apk apk update && apk upgrade && apk add --no-cache git

RUN mkdir /muffet
ADD . /muffet/
WORKDIR /muffet

RUN go get -d .
RUN go build -o main .

ENTRYPOINT ["/muffet/main"]
