FROM golang:1.5-onbuild

RUN ["go", "get", "-u", "gopkg.in/godo.v1/cmd/godo"]

CMD ["godo", "watch-server", "--watch"]
