FROM golang:1.12.4-alpine3.9

RUN apk add nmap
RUN apk add net-tools
RUN apk add curl

COPY src/ /

WORKDIR /
CMD go run *.go $MOCK
