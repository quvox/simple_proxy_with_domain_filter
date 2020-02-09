FROM golang:1.13.7-alpine3.11
MAINTAINER quvox

ENV HOMEDIR /opt/src
WORKDIR $HOMEDIR

RUN apk add make

COPY go/. $HOMEDIR/
RUN cd $HOMEDIR && make build && mkdir -p ../conf

COPY ./docker-entrypoint.sh /

EXPOSE 8080

ENTRYPOINT ["sh", "/docker-entrypoint.sh"]
