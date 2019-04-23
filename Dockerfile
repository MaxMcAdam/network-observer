FROM golang:1.12.4-apline3.9

RUN apt-get update -y
RUN apt-get install -y nmap
RUN apt-get install -y net-tools
RUN apt-get install -y curl

COPY src/ /

WORKDIR /
RUN go run /src/*.go
