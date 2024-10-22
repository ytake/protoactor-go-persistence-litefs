FROM alpine
ARG LITEFS_CONFIG=litefs.yml

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

ADD ./litefs/litefs.static-lease.yml /tmp/litefs.static-lease.yml
RUN cp /tmp/litefs.static-lease.yml /etc/litefs.static-lease.yml

RUN apk add bash fuse3 sqlite ca-certificates curl

ENTRYPOINT litefs mount
