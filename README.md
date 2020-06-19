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
