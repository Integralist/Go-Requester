#!/bin/bash

# allow `go get` for private repositories
git config --global url."git@github.com:".insteadOf "https://github.com/"

if [ "$1" = "ssh_setup" ]; then
  # start our ssh agent
  eval "$(ssh-agent -s)"

  # run expect to handle entering password for my mounted SSH key
  # /ssh.exp
  ssh-add /go/src/app/github_rsa

  # automate trusting github as a remote hote
  ssh -o StrictHostKeyChecking=no git@github.com uptime
  ssh -T git@github.com
fi

if [ ! -d "vendor/src" ]; then
  echo "Fetching dependencies described within manifest..."
  gb vendor restore
fi

if [ ! -d "/go/src/vendor" ]; then
  echo "Copying vendor packages into \$GOPATH..."
  cp -r /go/src/app/vendor/src/* /go/src/

  echo "Below is the container's \$GOPATH (with dependencies copied over)..."
  tree -L 3 /go/src
fi

echo "Start watching files..."
godo watch-server --watch
