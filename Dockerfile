FROM golang:alpine AS build
RUN apk add --no-cache gcc musl-dev
COPY . /tmp/src
WORKDIR /tmp/src
RUN mkdir -p /tmp/build
RUN go mod download
RUN go build -o /tmp/build/api

FROM alpine:latest
COPY --from=build /tmp/build/api /api
ENTRYPOINT ["/api"]
EXPOSE 8000
