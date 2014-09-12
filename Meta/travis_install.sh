#!/bin/sh -eux

# Values of $TRAVIS_GO_VERSION: go1.1 go1.3.1 tip

case ${TRAVIS_GO_VERSION#go} in
1.1|1.1.*)
  # No coverage support, avoid installing the tools needlessly
  ;;
*)
  go get code.google.com/p/go.tools/cmd/cover
  go get github.com/mattn/goveralls
  ;;
esac
