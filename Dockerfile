FROM golang:1.26.3@sha256:2981696eed011d747340d7252620932677929cce7d2d539602f56a8d7e9b660b AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
