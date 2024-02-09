emailsupport
============

This package contains auxiliary support information and routines for dealing
with email handling.

[![Build Status](https://api.travis-ci.org/philpennock/emailsupport.svg?branch=main)](https://travis-ci.org/philpennock/emailsupport)
[![Documentation](http://godoc.org/github.com/philpennock/emailsupport?status.svg)](http://godoc.org/github.com/philpennock/emailsupport)
[![Coverage Status](https://coveralls.io/repos/philpennock/emailsupport/badge.svg?branch=main)](https://coveralls.io/r/philpennock/emailsupport?branch=main)

At present, it only has some regular expressions which have been tested by
being in use for many years, in Perl, but have here been translated to
Golang's regexp library.  Other bits and pieces will creep in, as this package
acts as a ‘miscellaneous’ catch-all for anything Golang that's email-related.
As such, I'm not prepared to make API guarantees, so be sure to use dependency
management to track this repository.


Using
-----

This package follows normal Go package naming convention and is `go get`
compatible.

The package is documented using the native godoc system.
A public interface is available at
[godoc.org](http://godoc.org/github.com/philpennock/emailsupport).

The allowed syntax for email addresses changes between [RFC2821][]/[RFC2822][]
and their replacements, [RFC5321][]/[RFC5322][].
By default, the regular expressions employ the newer syntax definitions, but
you can build the library with a build-tag of `rfc2822` to use the definitions
supplied in [RFC2822][] instead of those from [RFC5321][].

This package uses [semantic versioning](https://semver.org/).  
Note that Go only supports the most recent two minor versions of the language;
for the purposes of semver, we do not consider it a breaking change to add a
dependency upon a language or standard library feature supported by all
currently-supported releases of Go.


Tools
-----

This is primarily a library package.
It does include two commands though.
This follows the standard Go idiom of using sub-directories of `./cmd` to hold the commands.
Thus you can use `go install ./cmd/...` to install them.

Or: `go install -v github.com/philpennock/emailsupport/cmd/...@latest`

 1. `email-regexp-emit`: just prints a regular expression for an email
    address.  The pattern uses `(?:  )` as a non-capturing group and is
    otherwise a simple Extended Regular Expression, so just about any modern
    regular expression library should be able to use it.

 2. `check-is-emailaddr`: can be given regexps on the command-line, or via an
    input file, and for each one reports success or failure.
    It exits true (0) if and only if every address given is fine.
    It exits 1 if some input is not an email address.
    It exists another non-zero value for problems in running.


Testing
-------

Run `go test`

[build-tag]: http://golang.org/pkg/go/build/#hdr-Build_Constraints
             "Build Constraints"
[RFC2821]: https://www.ietf.org/rfc/rfc2821.txt
           "Simple Mail Transfer Protocol"
[RFC2822]: https://www.ietf.org/rfc/rfc2822.txt
           "Internet Message Format"
[RFC5321]: https://www.ietf.org/rfc/rfc5321.txt
           "Simple Mail Transfer Protocol"
[RFC5322]: https://www.ietf.org/rfc/rfc5322.txt
           "Internet Message Format"
