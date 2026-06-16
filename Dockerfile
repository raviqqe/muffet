FROM golang:1.26.4@sha256:792443b89f65105abba56b9bd5e97f680a80074ac62fc844a584212f8c8102c3 AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
