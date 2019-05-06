FROM golang:1.12-alpine as builder
RUN apk add --no-cache git
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o showrss .


FROM alpine:3.9
MAINTAINER Fabien Foerster <fabienfoerster@gmail.com>
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/showrss /usr/bin/showrss
ENTRYPOINT ["showrss"]
