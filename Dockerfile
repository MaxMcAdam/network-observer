FROM ubuntu:latest

RUN apt-get update -y
RUN apt-get install -y nmap
RUN apt-get install -y net-tools
RUN apt-get install -y curl

COPY *.sh /
WORKDIR /
