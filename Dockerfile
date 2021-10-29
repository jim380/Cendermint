# FROM golang:1.15
FROM golang:alpine3.13 AS build-env

# Set up dependencies
ENV PACKAGES bash curl make git libc-dev gcc linux-headers eudev-dev python3

# ADD . /cendermint
WORKDIR /cendermint

COPY go.mod .
COPY go.sum .

COPY . .

RUN apk add --no-cache $PACKAGES && go build

FROM alpine:edge

RUN apk add --update ca-certificates

WORKDIR /cendermint

COPY --from=build-env /cendermint/Cendermint /usr/bin/Cendermint

CMD ["Cendermint"]