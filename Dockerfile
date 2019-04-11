FROM ubuntu:latest

COPY *.sh /
WORKDIR /
CMD /service.sh
