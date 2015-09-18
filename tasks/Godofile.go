package main

import . "gopkg.in/godo.v1"

func tasks(p *Project) {
	p.Task("watch-server", func(c *Context) error {
		// Run("gb build all")
		// return Run("./bin/requester ./src/page.yaml")
		// Avoiding because it doesn't trigger rebuild
		return Start("main.go ../page.yaml", M{"$in": "./src/requester"})
		// Either way I need to `docker stop <cid>`
	}).Watch("**/*.go")
}

func main() {
	Godo(tasks)
}
