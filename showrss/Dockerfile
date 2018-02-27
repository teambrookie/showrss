FROM alpine:3.4
MAINTAINER Fabien Foerster <fabienfoerster@gmail.com>
RUN apk add --no-cache ca-certificates
ADD showrss /usr/bin/showrss
ENTRYPOINT ["showrss"]
