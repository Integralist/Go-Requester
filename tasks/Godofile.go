package main

import . "gopkg.in/godo.v1"

func tasks(p *Project) {
	p.Task("run", func(c *Context) error {
		return Start(`go run requester.go`)
	}).Watch("**/*.go")
}

func main() {
	Godo(tasks)
}
