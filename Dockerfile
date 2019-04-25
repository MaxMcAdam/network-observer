FROM golang:1.12.4-alpine3.9

RUN apk add nmap net-tools curl mosquitto-clients

COPY src/ /
COPY scan.xml /
COPY servenv /

WORKDIR /
CMD source servenv
CMD go run *.go $COUCHDB_URL $WIOTP_ORG $WIOTP_DEVICE_TYPE $WIOTP_DEVICE_ID $WIOTP_DEVICE_TOKEN $MOCK
