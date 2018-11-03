FROM alpine
MAINTAINER Christopher Hein <me@chrishein.com>

RUN apk --no-cache add openssl musl-dev ca-certificates libc6-compat
COPY eniconfig-controller /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/eniconfig-controller"]
