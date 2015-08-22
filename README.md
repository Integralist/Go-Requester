<h1 align="center">Go-Requester</h1>

<p align="center">
  <img src="https://img.shields.io/badge/TODO-5%25-green.svg?style=flat-square">
</p>

<p align="center">
  <b>HTTP service</b> that accepts a collection of "components"<br>then fans-out requests and returns <b>aggregated content</b>
</p>

> Requires Go 1.5+

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

> Note: example config provided

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
- `go run slow-endpoint.go`
- `curl http://localhost:8080/` (better to check via a web browser)

### TODO

- Tests!

### Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
