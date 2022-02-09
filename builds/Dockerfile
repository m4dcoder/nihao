FROM alpine:latest

RUN mkdir -p /opt/nihao
COPY nihao /opt/nihao/

EXPOSE 6688/tcp
ENTRYPOINT ["/opt/nihao/nihao"]
