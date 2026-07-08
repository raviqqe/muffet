FROM golang:1.27rc2@sha256:2317c8e806fe884a0d5f3d80d596dcd1369a4c211b7eba8d507718245e9e5831 AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
