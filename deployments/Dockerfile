FROM golang:1.21.10-alpine AS build

RUN apk --no-cache add make

WORKDIR /post-service

COPY go.mod go.sum /post-service/
RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /

COPY --from=build /post-service/bin/app .
# COPY --from=build /post-service/configs configs

ENTRYPOINT [ "./app" ]