FROM alpine:3.4
MAINTAINER Fabien Foerster <fabienfoerster@gmail.com>
RUN apk add --no-cache ca-certificates
ENV GOOGLE_APPLICATION_CREDENTIALS /usr/var/googlecloud/key.json
ADD showrss /usr/bin/showrss
ADD showrss_service_key.json /usr/var/googlecloud/key.json

ENTRYPOINT ["showrss"]
