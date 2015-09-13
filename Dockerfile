FROM golang:1.5-onbuild

RUN ["go", "get", "-u", "gopkg.in/godo.v1/cmd/godo"]
RUN ["go", "get", "-u", "github.com/constabulary/gb/..."]
RUN ["gb", "vendor", "fetch", "gopkg.in/yaml.v2"]

CMD ["godo", "watch-server", "--watch"]
