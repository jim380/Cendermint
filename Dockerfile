FROM golang:1.18 AS build-env
ENV CGO_ENABLED=0
WORKDIR /cendermint
COPY . .
RUN go build

FROM alpine:3.15
WORKDIR /cendermint
COPY --from=build-env /cendermint/Cendermint /usr/bin/cendermint
CMD ["cendermint"]
