#!/bin/sh -eux

PATH="${GOPATH%%:*}/bin:$PATH"

: ${BUILDTAGS?'Need BUILDTAGS defined (may be empty)'}

case ${TRAVIS_GO_VERSION#go} in
1.1|1.1.*)
  # No coverage support
  go test -v -tags "$BUILDTAGS" ./...
  ;;
*)
  set +x ; : ${COVERALLS_TOKEN:?Need the coveralls.io token} ; set -x
  go test -v -tags "$BUILDTAGS" -covermode=count -coverprofile=coverage.out ./...
  set +x
  echo 'goveralls -coverprofile=coverage.out -service=travis-ci -repotoken [censored] (failure ignored)'
  goveralls -coverprofile=coverage.out -service=travis-ci -repotoken "$COVERALLS_TOKEN" || true
  set -x
  ;;
esac
