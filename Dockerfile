FROM golang:latest

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "tree", "expect"]
RUN ["go", "get", "-u", "gopkg.in/godo.v2/cmd/godo"]
RUN ["go", "get", "-u", "github.com/Masterminds/glide"]

ENV GO15VENDOREXPERIMENT=1

COPY ssh.exp /
RUN chmod +x /ssh.exp

COPY bootstrap.sh /
RUN chmod +x /bootstrap.sh

WORKDIR /go/src/app

EXPOSE 8080

ENTRYPOINT ["/bin/bash", "/bootstrap.sh"]
