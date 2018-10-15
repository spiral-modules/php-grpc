FROM php:7.2.10-alpine

RUN apk add make autoconf automake libc-dev gcc \
    && pecl install protobuf-3.6.1 \
    && echo "extension=protobuf.so" > /usr/local/etc/php/conf.d/protobuf.ini \
    && apk del make autoconf automake libc-dev gcc