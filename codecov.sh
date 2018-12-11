#!/bin/sh

if [ "$CODECOV_TOKEN" != "" ]; then
  go get github.com/pierrre/gotestcover
  gotestcover -coverprofile=coverage.txt ./...
  curl -s https://codecov.io/bash | bash -s --
else
  echo "please specify \$CODECOV_TOKEN"
  echo "https://github.com/cainus/codecov.io/blob/master/README.md#upload-repo-tokens"
fi
