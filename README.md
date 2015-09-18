<h1 align="center">Go-Requester</h1>

<p align="center">
  <img src="https://img.shields.io/badge/Completed-90%25-green.svg?style=flat-square">
</p>

<p align="center">
  <b>HTTP service</b> that accepts a collection of "components"<br>then fans-out requests and returns <b>aggregated content</b>
</p>

## Summary

- Components should be defined in a YAML page configuration file 
- Components are requested concurrently via goroutines
- Components can be marked as "mandatory" (if they fail, the request summary is set to "failure")

## Example Page Config

```yaml
components:
  - id: google
    url: http://google.com
  - id: integralist
    url: http://integralist.co.uk
    mandatory: true
  - id: not-found
    url: http://httpstat.us/404
  - id: timeout
    url: http://httpstat.us/408
  - id: server-error
    url: http://httpstat.us/500
  - id: service-unavailable
    url: http://httpstat.us/503
```

> Note: example config provided as part of this repo

## Example JSON Output

```json
{
  "summary": "success",
  "components": [
    {
      "id": "google",
      "status": 200,
      "body": "<!doctype html> ... </html>",
      "mandatory": false
    },
    {
      "id": "integralist",
      "status": 200,
      "body": "<!doctype html> ... </html>",
      "mandatory": true
    },
    {
      "id": "slooow",
      "status": 408,
      "body": "Get http://localhost:3000/pugs: net/http: request canceled (Client.Timeout exceeded while awaiting headers)",
      "mandatory": false
    }
    {
      "id": "not-found",
      "status": 404,
      "body": "404 Not Found",
      "mandatory": false
    }
  ]
}
```

> Note: the toplevel `summary` key's value will be `failure` if any mandatory components fail

## Run within Docker

- `docker build -t my-golang-app .`
- `docker run --rm --name go-tester -v "$PWD":/go/src/app -p 8080:8080 my-golang-app`

> To test: `docker run -it -v "$PWD":/go/src/app -p 8080:8080 my-golang-app /bin/bash`

## Build and run binary on host machine

The following only needs to be run once:

```bash
go get -u github.com/constabulary/gb/...
gb vendor fetch gopkg.in/yaml.v2
```

> Note: to get 'all' dependencies you can also use  
> gb vendor fetch github.com/<user>/<repo>

Every time you make a change to your code, run:

```bash
gb build all && bin/requester ./src/page.yaml
```

### Run application locally on host machine

- `go run src/requester/main.go src/page.yaml`
- `go run slow-endpoint.go` (see below for an example script)
- `curl http://localhost:8080/` (better to check via a web browser)

> Note: you can also use `godo watch-server --watch` to track changes and automatically re-run

## Slow HTTP Server

```go
package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5000 * time.Millisecond)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}

// Example -> http://localhost:3000/pugs
```

## Curl Timings

- `curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/`

> Note: I've included a `curl-format.txt` file within the repo

## Dependencies

I use http://getgb.io/ for handling dependencies. When using `gb vendor fetch <pkg>` it'll place dependencies into a `vendor` directory for you and thus allow `gb build all` to include them within your binary. So you gain a project specific workspace without affecting your global `$GOPATH`.

> Note: typically you'd execute `go vendor fetch github.com/integralist/go-requester` to pull in all dependencies

We `.gitignore` the `vendor/src` directory but we commit the `vendor/manifest` file (which acts as a dependency lockfile) and then when pulling this repo for the first time you'd need to execute `gb vendor restore`

## Compilation

Use http://getgb.io/ again, this time `go build all`

An alternative is to use [Gox](https://github.com/mitchellh/gox):

- `go get github.com/mitchellh/gox`
- `gox`

But I've not yet used it alongside `gb` so I'm not sure if there are any nuances to the setup.

## TODO

- Check use of `gb` to build different OS and ARCH binaries and include notes in README
- See if gox works alongside gb
- Switch from string checking 'timeout' error to using type assertion
- Add logic for loading page config remotely
- Dynamically change port number when run as binary
- Tests!

## Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
