#!/bin/bash

# allow `go get` for private repositories
git config --global url."git@github.com:".insteadOf "https://github.com/"

if [ "$1" = "ssh_setup" ]; then
  # start our ssh agent
  eval "$(ssh-agent -s)"

  # run expect to handle entering password for my mounted SSH key
  # /ssh.exp
  ssh-add /.ssh/github_rsa

  # automate trusting github as a remote hote
  ssh -o StrictHostKeyChecking=no git@github.com uptime
  ssh -T git@github.com
fi

if [ ! -d "vendor" ]; then
  echo "Fetching dependencies described within lock file..."
  glide install
fi

echo "Start watching files..."
godo server --watch
