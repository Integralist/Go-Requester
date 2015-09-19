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

If you're using Docker-Machine, then executing the following command will return the results of the running application:

```bash
curl $(docker-machine ip <name_of_vm>):8080
```

To test what's happening inside the container then execute the following command:

```bash
docker run -it -v "$PWD":/go/src/app -p 8080:8080 my-golang-app /bin/bash
```

### Private Repos

If you need to use a private repo then I've set-up the Docker container to work with them, but it does require you to pick your poison:

- Run the container interactively (`-it`) and manually enter SSH private key passphrase

OR

- Run the container interactively (`-it`) and save your SSH private key passphrase into an `expect` script

The former is safer and so it's the default option.

The latter requires you to modify the `ssh.exp` file to include your passphrase and also to uncomment `ssh-add /go/src/app/github_rsa` from `bootstrap.sh` and put back in `/ssh.exp`. But it's your responsibility to make sure you don't commit that passphrase.

> Yes the latter takes a few more steps, but there was no way I was going to automate that for you and make it easy for you to shoot yourself in the foot!

Regardless of which option you choose, you'll need to modify the Docker run command slightly (to use interactive mode `-it`). I'm not sure why though, if anyone has any ideas then please let me know. 

So the command you need to execute, if you're pulling private repositories, is:

```
docker run -it -v "$HOME/.ssh/github_rsa":/go/src/app/github_rsa -v "$PWD":/go/src/app -p 8080:8080 my-golang-app
```

> Note: I initially added both `github_rsa` and `id_rsa` to the `.gitignore` file and then switched to `*_rsa` to try and prevent as many mistakes as possible (i.e. avoiding committing your private key to a potentially public repo); so please do be careful if your private key uses a different name!

### In Summary...

No, I'm not specifying any private repo dependencies:

```bash
docker run --rm -v "$PWD":/go/src/app -p 8080:8080 my-golang-app
```

Yes, I have some dependencies coming from private repositories:

```bash
docker run -it -v "$HOME/.ssh/github_rsa":/go/src/app/github_rsa -v "$PWD":/go/src/app -p 8080:8080 my-golang-app "ssh_setup"
```

- run interactively (`-it` mode)
- mount your ssh private key
- specify `ssh_setup`

## Build and run binary on host machine

The following command only needs to be run once (it downloads the `gb` tool):

```bash
go get -u github.com/constabulary/gb/...
```

Once you have `gb`, download the dependencies (specified within the `vendor/manifest` file) using: 

```bash
gb vendor restore
```

Every time you make a change to your code, run:

```bash
gb build all && bin/requester ./src/page.yaml
```

### Run application locally on host machine

- `go run src/requester/main.go src/page.yaml`
- `curl http://localhost:8080/` (better to check via a web browser)

> Note: you can also use `godo watch-server --watch` to track changes and automatically re-run

If you want to see how latency/slow responses effect the application then try also running: https://github.com/Integralist/go-slow-http-server which was specifically designed to be used alongside go-requester for testing

## Curl Timings

- `curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8080/`

> Note: I've included a `curl-format.txt` file within the repo

## Dependencies

I use http://getgb.io/ for handling dependencies. When using `gb vendor fetch <pkg>` it'll place dependencies into a `vendor` directory for you and thus allow `gb build all` to include them within your binary. So you gain a project specific workspace without affecting your global `$GOPATH`.

In order to not have a large repo (consisting of many dependency files), we `.gitignore` the `vendor/src` directory but we still commit the `vendor/manifest` file (which acts as a dependency lockfile). This means when pulling the repo for the first time you'd need to execute `gb vendor restore`.

## Compilation

Use http://getgb.io/ again, this time `go build all`

An alternative is to use [Gox](https://github.com/mitchellh/gox):

- `go get github.com/mitchellh/gox`
- `gox`

But I've not yet used it alongside `gb` so I'm not sure if there are any nuances to the setup.

## TODO

- Check use of `gb` to build different OS and ARCH binaries and include notes in README
- See if gox works alongside gb
- Add logic for loading page config remotely
- Dynamically change port number when run as binary
- Tests!

## Licence

[The MIT License (MIT)](http://opensource.org/licenses/MIT)

Copyright (c) 2015 [Mark McDonnell](http://twitter.com/integralist)
