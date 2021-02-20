FROM golang
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go get .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
