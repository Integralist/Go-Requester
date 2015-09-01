package main

import . "gopkg.in/godo.v1"

func tasks(p *Project) {
	p.Task("watch-server", func(c *Context) error {
		return Start(`requester.go`)
	}).Watch("**/*.go")
}

func main() {
	Godo(tasks)
}
