<h1 align="center">Go-Requester</h1>

<p align="center">
  <img src="https://img.shields.io/badge/TODO-5%25-green.svg?style=flat-square">
</p>

<p align="center">
  <b>HTTP service</b> that accepts a collection of "components"<br>then fans-out requests and returns <b>aggregated content</b>
</p>

> Requires Go 1.5+

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
      "body": "<doctype! html> ... </html>",
      "mandatory": false
    },
    {
      "id": "integralist",
      "status": 200,
      "body": "<doctype! html> ... </html>",
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

## Setup example

- `go run requester.go`
- `go run slow-endpoint.go` (see below for an example script)
- `curl http://localhost:8080/` (better to check via a web browser)

### Slow HTTP Server

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

## Compilation

I recommend using [Gox](https://github.com/mitchellh/gox).

- `go get github.com/mitchellh/gox`
- `gox`

## TODO

- Tests!

## Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
