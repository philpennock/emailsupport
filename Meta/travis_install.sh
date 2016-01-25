#!/bin/sh -eux

# Values of $TRAVIS_GO_VERSION are of general form: go1.1 go1.3.1 tip
# Precise values determined by `.travis.yml`

case ${TRAVIS_GO_VERSION#go} in
1.1|1.1.*)
  # No coverage support, avoid installing the tools needlessly
  ;;
1.2|1.2.*|1.3|1.3.*|1.4|1.4.*)
  go get golang.org/x/tools/cmd/cover
  go get github.com/mattn/goveralls
  ;;
*)
  # Golang 1.5 and later, cover is in the standard repo
  go get github.com/mattn/goveralls
  ;;
esac
