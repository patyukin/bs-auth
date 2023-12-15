FROM alpine:3.18

RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cacher/apk/*

ADD https://github.com/pressly/goose/releases/download/v3.14.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

WORKDIR /root

COPY .env.local .env
COPY ./migrations/*.sql migrations/
COPY migrations.sh .

RUN chmod +x migrations.sh

ENTRYPOINT ["bash", "migrations.sh"]
