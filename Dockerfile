FROM golang:1.26.5@sha256:079e59808d2d252516e27e3f3a9c003740dee7f75e55aa71528766d52bcfc16a AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
