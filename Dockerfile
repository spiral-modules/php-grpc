FROM php:7.2.11-cli-alpine

ENV PROTOBUF_VERSION 3.6.1

RUN apk add make autoconf automake libc-dev gcc \
    && pecl install protobuf-$PROTOBUF_VERSION \
    && echo "extension=protobuf.so" > /usr/local/etc/php/conf.d/protobuf.ini \
    && apk del make autoconf automake libc-dev gcc

RUN cd /tmp \
    && apk add libc6-compat \
    && wget https://github.com/protocolbuffers/protobuf/releases/download/v$PROTOBUF_VERSION/protoc-$PROTOBUF_VERSION-linux-x86_64.zip \
    && unzip -p protoc-$PROTOBUF_VERSION-linux-x86_64.zip bin/protoc > /usr/bin/protoc \
    && chmod +x /usr/bin/protoc \
    && rm protoc-$PROTOBUF_VERSION-linux-x86_64.zip

ENTRYPOINT [""]