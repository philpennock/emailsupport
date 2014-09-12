// © Phil Pennock 2013.  See LICENSE file for licensing.

/*
Package emailsupport contains some more esoteric routines useful for parsing
and handling email.  For instance, some baroque regular expressions.

No APIs are stable.  Make sure you handle your dependencies accordingly.


REGULAR EXPRESSIONS

The package creates a number of exported regular expression objects,
initialised at init time, pulling up the creation cost to program start.  Per
the documentation for the `regexp` package, “A Regexp is safe for concurrent
use by multiple goroutines.”

Each regular expression comes in two forms, `Foo` and `FooUnanchored`.  The
short name is anchored to the start and end of the pattern namespace, so that
by default if you match an variable against `EmailAddress`, you are confirming
that the contents of the variable are an email address, not accidentally only
matching that there's something address-like somewhere within.  The longer name
(I need something briefer which is still clear) provides a pattern which can
be used to find a regular expression elsewhere.

Each regular expression is available as a pattern string in a form with a
`Txt` prefix, which can be used to build larger regular expressions.  The
pattern is wrapped with `(?:...)` to be a non-capturing group which can be
qualified or otherwise treated as a single unit.  No regular expression has
any capturing groups, letting the caller manage capturing indices.

Thus for any `Foo`, this package provides:

 * TxtFoo: the regular expression as a text string, safe for embedding
 * Foo: a regexp object which is anchored to start and end
 * FooUnanchored: a regexp object which is not anchored

The list includes:

 * `EmailAddress`: an RFC5321 email address, the part within angle-brackets or
   as is used in SMTP.  This is _not_ an RFC5322 message header email email
   address, with display forms and comments.
 * `EmailDomain`: a domain which can be used in email addresses; this is the
   base specification form and does not handle internationalisation, though this
   regexp should be correct to apply against punycode-encoded domains.  This does
   handle embedded IPv4 and IPv6 address literals, but not the
   General-address-literal which is a grammar hook for future extension (because
   that inherently can't be handled until defined).
 * `EmailLHS`: the Left-Hand-Side (or "local part") of an email address; this
   handles unquoted and quoted forms.
 * `EmailAddressOrUnqualified`: either an address or a LHS, this is a form often
   used in mail configuration files where a domain is implicit.

 * `IPv4Address`, `IPv6Address`: an IPv4 or IPv6 address
 * `IPv4Netblock`, `IPv6Netblock`, IPNetblock: a netblock in CIDR prefix/len notation (used
   for source ACLs)
 * `IPv4Octet`: a number 0 to 255

The IPv6 address regexp is taken from RFC3986 (the one which gets it right) and
is a careful copy/paste and edit of a version which has been used and gradually
debugged for years, including in a tool I released called `emit_ipv6_regexp`.

The patterns used in `EmailLHS` (and thus also in items which include an email
left-hand-side) can be in one of two forms, and selecting between them is a
compile-time decision.  The rules can be either those from RFC2822 or those from
RFC5321.  By default, those from RFC5321 are used.  Build with a `rfc2822`
build-tag to get the older definitions.  If a future RFC changes the rules
again, then the default patterns in this package may change; the build-tag
`rfc5321` is currently unused, but is reserved for the future to force
selecting the rules which are now current.
It is safe (harmless) to supply that build-tag now (but not together with
`rfc2822`).
*/
package emailsupport

// I know full well that markdown isn't handled in godoc; I choose to use the
// notation anyway, in the hopes of future tooling that provide for better
// markup control.
