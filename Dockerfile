FROM golang:1.5

RUN ["go", "get", "-u", "gopkg.in/godo.v1/cmd/godo"]
RUN ["go", "get", "-u", "github.com/constabulary/gb/..."]

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

CMD ["gb", "vendor", "fetch", "gopkg.in/yaml.v2", "&&", "godo", "watch-server", "--watch"]
