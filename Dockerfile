FROM alpine:3.4
MAINTAINER Fabien Foerster <fabienfoerster@gmail.com>
ADD showrss /usr/bin/showrss
ENTRYPOINT ["showrss"]
