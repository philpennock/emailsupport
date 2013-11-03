// © Phil Pennock 2013.  See LICENSE file for licensing.

/*
This package contains some more esoteric routines useful for parsing
and handling email.  For instance, some baroque regular expressions.

No APIs are stable.  Make sure you handle your dependencies accordingly.
*/
package emailsupport

import (
	"regexp"
	"strings"
	"unicode"
)

// These regular expressions come from my "update_mail_cdbs" tool, are
// copy&pasted across and mostly-only adjusted for Perl->Go syntax.

// reduce chance of typos in start/end by letting parser detect typos
const (
	start = `\A`
	end   = `\z`
)

const (
	txtIpv4Octet    = `(?:25[0-5]|(?:[1-9]|1[0-9]|2[0-4])?[0-9])`
	txtIpv4Address  = `(?:(?:` + txtIpv4Octet + `\.){3}` + txtIpv4Octet + `)`
	txtIpv4Netblock = `(?:` + txtIpv4Address + `/(?:[12]?[0-9]|3[0-2]))`
)

var (
	Ipv4OctetUnanchored    = regexp.MustCompile(txtIpv4Octet)
	Ipv4Octet              = regexp.MustCompile(start + txtIpv4Octet + end)
	Ipv4AddressUnanchored  = regexp.MustCompile(txtIpv4Address)
	Ipv4Address            = regexp.MustCompile(start + txtIpv4Address + end)
	Ipv4NetblockUnanchored = regexp.MustCompile(txtIpv4Netblock)
	Ipv4Netblock           = regexp.MustCompile(start + txtIpv4Netblock + end)
)

/*
  RFC 3986 states:
        IPv6address =                            6( h16 ":" ) ls32
                    /                       "::" 5( h16 ":" ) ls32
                    / [               h16 ] "::" 4( h16 ":" ) ls32
                    / [ *1( h16 ":" ) h16 ] "::" 3( h16 ":" ) ls32
                    / [ *2( h16 ":" ) h16 ] "::" 2( h16 ":" ) ls32
                    / [ *3( h16 ":" ) h16 ] "::"    h16 ":"   ls32
                    / [ *4( h16 ":" ) h16 ] "::"              ls32
                    / [ *5( h16 ":" ) h16 ] "::"              h16
                    / [ *6( h16 ":" ) h16 ] "::"

        ls32        = ( h16 ":" h16 ) / IPv4address
                    ; least-significant 32 bits of address

        h16         = 1*4HEXDIG
                    ; 16 bits of address represented in hexadecimal

  Note that we need to allow:
    1:2:3:4:5:6:7:8   -- 7 colons
    ::2:3:4:5:6:7:8   -- 8 colons, skips leading 0
    1::3:4:5:6:7:8    -- 7 colons again; or fewer
    1:2:3:4:5:6::8    -- 7 or fewer
    1:2:3:4:5:6:7::   -- 8 colons, skips trailing 0
  so there can be 8 colons only if two are doubled and are an affix.
  Otherwise there's always 7 colons at most.

  RFC 4291: IPv6 Addressing Architecture
    The use of "::" indicates one or more groups of 16 bits of zeros.
    The "::" can only appear once in an address.  The "::" can also be
    used to compress leading or trailing zeros in an address.
  That's "1 or more", not "2 or more", so in effect when it's an affix
  there's a degenerate case where a colon (:) just replaces a zero (0).

  Bugs encountered after original writing (in Perl):
   * Was missing two /o optimisations
   * Had extra | at end of the last line, permitting empty string to match
*/

// this would be easier if the Go constructor supported extended mode

const (
	txtIpv6H16  = `(?:[0-9a-fA-F]{1,4})`
	txtIpv6LS32 = `(?:(?:` + txtIpv6H16 + `:` + txtIpv6H16 + `)|` + txtIpv4Address + `)`
)

var (
	txtIpv6Address = strings.Join([]string{`(?:`,
		`(?:(?:(?:`, txtIpv6H16, `:){6})`, txtIpv6LS32, `)|`,
		`(?:(?:::(?:`, txtIpv6H16, `:){5})`, txtIpv6LS32, `)|`,
		`(?:(?:(?:`, txtIpv6H16, `)?::(?:`, txtIpv6H16, `:){4})`, txtIpv6LS32, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,1}`, txtIpv6H16, `)?::(?:`, txtIpv6H16, `:){3})`, txtIpv6LS32, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,2}`, txtIpv6H16, `)?::(?:`, txtIpv6H16, `:){2})`, txtIpv6LS32, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,3}`, txtIpv6H16, `)?::`, txtIpv6H16, `:)`, txtIpv6LS32, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,4}`, txtIpv6H16, `)?::)`, txtIpv6LS32, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,5}`, txtIpv6H16, `)?::)`, txtIpv6H16, `)|`,
		`(?:(?:(?:(?:`, txtIpv6H16, `:){0,6}`, txtIpv6H16, `)?::))`,
		`)`}, "")
	txtIpv6Netblock = `(?:` + txtIpv6Address + `/(?:[1-9]?[0-9]|1[01][0-9]|12[0-8]))`
	txtIpNetblock   = `(?:` + txtIpv4Netblock + `|` + txtIpv6Netblock + `)`
)

var (
	Ipv6AddressUnanchored  = regexp.MustCompile(txtIpv6Address)
	Ipv6Address            = regexp.MustCompile(start + txtIpv6Address + end)
	Ipv6NetblockUnanchored = regexp.MustCompile(txtIpv6Netblock)
	Ipv6Netblock           = regexp.MustCompile(start + txtIpv6Netblock + end)

	// these don't handle scoped addresses, but SMTP doesn't permit them
)

// See RFC 2821 (not all comments are terminology from there)
// Also atext and qtext from RFC 2822
// NB: this relies upon ASCII ordering; should not pose a problem

func deExtend(extended string) string {
	inLines := strings.Split(extended, "\n")
	outLines := make([]string, 0, len(inLines))
	for i := range inLines {
		c := strings.TrimSpace(inLines[i])
		if strings.HasPrefix(c, "#") {
			continue
		}
		outLines = append(outLines, strings.Map(
			func(r rune) rune {
				if unicode.IsSpace(r) {
					return -1
				} else {
					return r
				}
			},
			c))
	}
	return strings.Join(outLines, "")
}

const (
	txtAText = `[A-Za-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]`
	txtQText = `[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]`
)

var txtEmailLHS string = deExtend(`
	 # Local-part
	 (?:
	  (?:
		# Dot-string
		(?:` + txtAText + `)+ (?: \. ` + txtAText + `+)*
	  ) | (?:
		# Quoted-string
		" (?: \s* (?:
		 (?: ` + txtQText + `+ ) |
		 (?: \\ [\x01-\x09\x0b\x0c\x0e-\x7f] )
		) )* "
	  )
	 )`)

var txtEmailDomain string = deExtend(`
	 # Domain
	 (?:
		(?:
		 # regular domain
		 (?:[A-Za-z0-9] (?: [A-Za-z0-9-]*[A-Za-z0-9] )?)
		 (?: \. [A-Za-z0-9] (?: [A-Za-z0-9-]*[A-Za-z0-9] )?)+
		) | (?:
			 # address-literals
			 \[
			   (?: ` + txtIpv4Address + ` | (?: [Ii][Pp][vV]6: ` + txtIpv6Address + ` ) )
			   # NOTSUPP: General-address-literal
			   # G-A-L is a hook for future literal addresses
			   # and only specifies tag:content
			 \]
		)
	 )`)

var txtEmailAddress = `(?:(?:` + txtEmailLHS + `)@(?:` + txtEmailDomain + `))`

var txtEmailAddressOrUnqualified = `(?:` + txtEmailLHS + `(?:@` + txtEmailDomain + `)?)`

var (
	EmailLHSUnanchored                  = regexp.MustCompile(txtEmailLHS)
	EmailLHS                            = regexp.MustCompile(start + txtEmailLHS + end)
	EmailDomainUnanchored               = regexp.MustCompile(txtEmailDomain)
	EmailDomain                         = regexp.MustCompile(start + txtEmailDomain + end)
	EmailAddressUnanchored              = regexp.MustCompile(txtEmailAddress)
	EmailAddress                        = regexp.MustCompile(start + txtEmailAddress + end)
	EmailAddressOrUnqualifiedUnanchored = regexp.MustCompile(txtEmailAddressOrUnqualified)
	EmailAddressOrUnqualified           = regexp.MustCompile(start + txtEmailAddressOrUnqualified + end)
)
