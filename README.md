emailsupport
============

This package contains auxiliary support information and routines for dealing
with email handling.

At present, it only has some regular expressions which have been tested by
being in use for many years, in Perl, but have here been translated to
Golang's regexp library.  Other bits and pieces will creep in, as this package
acts as a ‘miscellaneous’ catch-all for anything Golang that's email-related.
As such, I'm not prepared to make API guarantees, so be sure to use dependency
management to track this repository.


Using
-----

This package follows normal Go package naming convention and is `go get`
compatible.  Additionally, the [godep](https://github.com/kr/godep) command is
used to track dependencies (currently none), in the file `Godeps`.

The package is documented using the native godoc system.

<!-- TODO: include a godoc.org link once repository is public -->


Testing
-------

Run `go test`
