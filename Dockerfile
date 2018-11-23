FROM golang
ADD . /go/src/github.com/raviqqe/muffet
RUN CGO_ENABLED=0 GOOS=linux go get /go/src/github.com/raviqqe/muffet

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
