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

> Note:  
> example config provided as part of this repo  
> `./config/page.yaml`

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

## Installing Vendored Dependencies

You'll need [Glide](https://github.com/Masterminds/glide) installed:

```bash
go get github.com/Masterminds/glide
```

If using Go `1.6` you're OK, otherwise you'll need to set the vendor experiment:

```bash
export GO15VENDOREXPERIMENT=1
```

To install the dependencies found in the `Glide.lock` file:

```bash
glide install
```

## Run Tests

```bash
go test $(glide novendor)
```

## Build Application

```bash
go build
```

Once the application is built and installed into your binary path, you can run it:

```bash
go-requester <path/to/your/page/config>
```

## Running Locally

There are two ways to run this application locally:

1. Host Machine
2. Docker

### Host Machine

```bash
go run main.go ./config/page.yaml
```

> Alternative: `godo server --watch`

Once the application is running, view: `http://localhost:8080/`

If you want to see how latency/slow responses effect the application, then also try running: https://github.com/Integralist/go-slow-http-server which was specifically designed to be used alongside go-requester for testing purposes.

### Docker

```bash
docker build -t gorequester .
docker run --rm --name gr -v "$PWD":/go/src/app -p 5000:8080 gorequester
```

Once the application is running, view: `http://<docker-machine-ip>:5000/`

To debug the container:

```bash
docker exec -it <container_id> /bin/bash
docker run -it --entrypoint /bin/bash -v "$PWD":/go/src/app -p 5000:8080 gorequester
```

If you specify a private repo as a dependency, then you'll need to pass in your SSH credentials:

```bash
docker run \
  -it \
  -v "$PWD":/go/src/app \
  -v "$HOME/.ssh/github_bbc_rsa":/.ssh/github_rsa \
  -p 5000:8080 \
  gorequester ssh_setup
```

Notice we've passed `ssh_setup`. Our bootstrap script (which runs inside the Docker container) will identify when this argument is provided and subsequently load up the SSH agent within the container. 

Because we've added the `-it` flag we're able to manually type in our SSH private key's passphrase (you *do* use a passphrase don't you?).

#### Expect

Having to type in your SSH passphrase can be annoying an restrict your ability to automate. Although I personally wouldn't advise it, you *can* make a change to the bootstrap script and have it use the `ssh.exp` script included in this repository.

Using this Expect script will mean you don't have to manually enter the passphrase for your SSH private key. The MASSIVE downside to this process is that you'll need to edit the `ssh.exp` script to include your passphrase. This means your secret passphrase has now been codified into the file.

If your laptop is compromised or you commit the change accidentally, then you're in serious trouble. So really... don't use it unless you're really sure you know what you're doing.

You'll also need to comment out the line `ssh-add /.ssh/github_rsa` from the bootstrap script and replace it with `/ssh.exp`.

## Curl Timings

We can verify the performance of this application using curl timings, like so:

```bash
curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/
```

> Note: I've included a `curl-format.txt` file within the repo

## TODO

- Refactor code so that certain aspects are loaded from other packages
- Add logic for loading page config remotely
- Dynamically change port number when run as binary
- Tests!

## Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
