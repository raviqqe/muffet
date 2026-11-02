FROM golang:1.27.3@sha256:2920bb124f1a2e12a1d39f0cab756ae0068441ee1da7655dec7f5e374944095b AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM scratch
COPY --from=build /go/bin/muffet /muffet
ENTRYPOINT ["/muffet"]
