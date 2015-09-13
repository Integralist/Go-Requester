package main

import . "gopkg.in/godo.v1"

func tasks(p *Project) {
	p.Task("watch-server", func(c *Context) error {
		return Start("main.go ../page.yaml", M{"$in": "./src/requester"})
	}).Watch("**/*.go")
}

func main() {
	Godo(tasks)
}
