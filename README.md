<h1 align="center">Go-Requester</h1>

<p align="center">
  <img src="https://img.shields.io/badge/status-90%25-green.svg?style=flat-square">
</p>

<p align="center">
  <b>HTTP service</b> that accepts a collection of "components"<br>then fans-out requests and returns <b>aggregated content</b>
</p>

> Requires Go 1.5+

## Example JSON Output

```json
{
  "status": "success",
  "components": [
    {
      "id": "google",
      "status": 200,
      "body": "<doctype! html> ... </html>"
    },
    {
      "id": "integralist",
      "status": 200,
      "body": "<doctype! html> ... </html>"
    },
    {
      "id": "slooow",
      "status": 408,
      "body": "Get http://localhost:3000/pugs: net/http: request canceled (Client.Timeout exceeded while awaiting headers)"
    }
    {
      "id": "not-found",
      "status": 404,
      "body": "404 Not Found"
    }
  ]
}
```

## Setup example

- `go run requester.go`
- `go run slow-endpoint.go`
- `curl http://localhost:8080/` (better to check via a web browser)

## TODO

- Add mandatory key
- Ensure overall status is set to fail if any mandatory components fail

### Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
