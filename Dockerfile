FROM golang:1.26.5@sha256:d52df9c279840adf958d017ebb275651ed8338b953d39817bc3633a2e6b1bbcc AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
