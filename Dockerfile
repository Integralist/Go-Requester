FROM golang:1.5

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "tree"]
RUN ["go", "get", "-u", "gopkg.in/godo.v1/cmd/godo"]
RUN ["go", "get", "-u", "github.com/constabulary/gb/..."]

COPY bootstrap.sh /
RUN chmod +x /bootstrap.sh

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["/bootstrap.sh"]
