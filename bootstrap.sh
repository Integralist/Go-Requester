#!/bin/bash

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
