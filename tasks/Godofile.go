package main

import (
	"fmt"
	"os"

	do "gopkg.in/godo.v2"
)

func tasks(p *do.Project) {
	if pwd, err := os.Getwd(); err == nil {
		do.Env = fmt.Sprintf("GOPATH=%s/vendor::$GOPATH", pwd)
	}

	p.Task("server", nil, func(c *do.Context) {
		c.Start("main.go ./config/page.yaml", do.M{"$in": "./"})
	}).Src("**/*.go")
}

func main() {
	do.Godo(tasks)
}
