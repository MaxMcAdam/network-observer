FROM golang:1.12.4-alpine3.9

RUN apk add nmap net-tools curl mosquitto-clients

COPY src/ /

COPY src/service-vars.json /

WORKDIR /
CMD go run *.go
