// © Phil Pennock 2013.  See LICENSE file for licensing.

package emailsupport

import (
	"regexp"
	"testing"
)

type boolPatternMatch struct {
	text        string
	shouldMatch bool
}

func iterateBoolPatternMatch(
	t *testing.T,
	pattern *regexp.Regexp,
	label string,
	list []boolPatternMatch,
) {
	for _, item := range list {
		if pattern.MatchString(item.text) != item.shouldMatch {
			var failType string
			if item.shouldMatch {
				failType = "does not match (but should)"
			} else {
				failType = "matches (but should not)"
			}
			t.Errorf("Regexp %s against %q %s", label, item.text, failType)
		}
	}
}

func TestIPv4Octets(t *testing.T) {
	iterateBoolPatternMatch(t, IPv4Octet, "IPv4Octet", []boolPatternMatch{
		{"0", true},
		{"1", true},
		{"9", true},
		{"10", true},
		{"25", true},
		{"26", true},
		{"99", true},
		{"100", true},
		{"101", true},
		{"156", true},
		{"199", true},
		{"200", true},
		{"201", true},
		{"240", true},
		{"245", true},
		{"246", true},
		{"249", true},
		{"250", true},
		{"251", true},
		{"252", true},
		{"253", true},
		{"254", true},
		{"255", true},
		{"256", false},
		{"260", false},
		{"1.1", false},
		{"-1", false},
		{"-255", false},
		{" 1 ", false},
	})
}

func TestIPv4Addresses(t *testing.T) {
	iterateBoolPatternMatch(t, IPv4Address, "IPv4Address", []boolPatternMatch{
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"0.0.0.0.0", false},
		{"192.0.2.255", true},
		{"192.0.256.250", false},
		{" 192.0.2.255", false},
		{"192.0.2.255.", false},
		{"192.168.1.2", true},
		{"...", false},
		{"192:0:2:2", false},
	})
}

func TestIPv4Netblocks(t *testing.T) {
	iterateBoolPatternMatch(t, IPv4Netblock, "IPv4Netblock", []boolPatternMatch{
		{"0.0.0.0/0", true},
		{"127.0.0.0/8", true},
		{"192.0.2.0/24", true},
		{"192.0.2.0/30", true},
		{"192.0.2.0/31", true},
		{"192.0.2.0/32", true},
		{"192.0.2.0/33", false},
		{"192.0.2.0/300", false},
		{"192.0.2.0", false},
		{"192.0.2.0/30 ", false},
		{"192.0.2.0/30/", false},
		{"192.0.2.0/30.", false},
	})
}

// these are the tests from my emit_ipv6_regexp tool
func TestIPv6AddressesFromEmitTester(t *testing.T) {
	iterateBoolPatternMatch(t, IPv6Address, "IPv6Address", []boolPatternMatch{
		{"::", true},
		{"::1", true},
		{"fe02::1", true},
		{"::ffff:192.0.2.1", true},
		{"2001:DB8::42", true},
		{"2001:db8::42", true},
		{"2001:DB8:1234:5678:90ab:cdef:0123:4567", true},
		{"2001:DB8:1234:5678:90ab:cdef:0123::", true},
		{"2001:DB8:1234:5678:90ab:cdef::0123", true},
		{"2001:DB8:1234:5678:90ab:cdef:192.0.2.1", true},
		{"2001:DB8:1234:5678:90ab:cdef:192.0.2.1", true},
		{"127.0.0.1", false},
		{"", false},
		{" ", false},
		{"192.0.2.1", false},
		{"2001", false},
		{"2001:DB8", false},
		{"2001:DB8:", false},
		{"2001:DB8::42::1", false},
		{"2001:DB8:1234:5678:90ab:cdef:g123:4567", false},
		{"2001:DB8:1234:5678:90ab:cdef:0123:4567:89", false},
		{"2001:DB8:1234:5678:90ab:cdef:0123", false},
	})
}

func TestIPv6Netblocks(t *testing.T) {
	iterateBoolPatternMatch(t, IPv6Netblock, "IPv6Netblock", []boolPatternMatch{
		{"::/0", true},
		{"fe02::/8", true},
		{"fe02::/08", false},
		{"fe02::/10", true},
		{"fe02::/16", true},
		{"2001:DB8:1234:5678::/64", true},
		{"2001:DB8:1234:5678::/127", true},
		{"2001:DB8:1234:5678::/128", true},
		{"2001:DB8:1234:5678::/129", false},
	})
}

func TestEmailLHS(t *testing.T) {
	iterateBoolPatternMatch(t, EmailLHS, "EmailLHS", []boolPatternMatch{
		{`john`, true},
		{`john.doe`, true},
		{`John.Doe`, true},
		{`alpha-beta`, true},
		{`john+topic`, true},
		{`""`, true},
		{`"john"`, true},
		{`"john doe"`, true},
		{`a~` + "`" + `*&^%$#!_-={|}'/?b`, true},
		{`#`, true},
		{`"X'); DROP TABLE domains; DROP TABLE passwords; --"`, true},
		{`"<script>alert('Boo!')</script>"`, true},
	})
}

func TestEmailDomain(t *testing.T) {
	iterateBoolPatternMatch(t, EmailDomain, "EmailDomain", []boolPatternMatch{
		{"example.org", true},
		{"example.org.", false},
		{".org", false},
		{"a-b.example", true},
		{"a--b.example", true},    // not valid to _register_ as a domain, but valid in SMTP grammar
		{"xn--4bi.example", true}, // xn--4bi = ✉ (ENVELOPE); xn-- being why -- is valid in a domain
		{"a-b", false},
		{"xn--4bi", false},
		{"", false},
		{".", false},
		{"192.0.2.1", true},   // is within a TLD 1, not for routing to an IP address
		{"[192.0.2.1]", true}, // routing to an IP address
		{"2001:db8::42", false},
		{"[2001:db8::42]", false},
		{"[ipv6:2001:db8::42]", true},
		{"[IPv6:2001:db8::42]", true},
	})
}

func TestEmailAddress(t *testing.T) {
	iterateBoolPatternMatch(t, EmailAddress, "EmailAddress", []boolPatternMatch{
		{`john@example.org`, true},
		{`john.doe@example.org`, true},
		{`sample-list@list.example.org`, true},
		{`john+foo@example.org`, true},
		{`<john.doe@example.org>`, false},
		{`john@our-subdomain`, false},
		{`john@our-subdomain.`, false},
		{`john@our-subdomain.example`, true},
		{`deliver@xn--4bi.example`, true},
		{`john@[IPv6:2001:db8::42]`, true},
		{`john@[192.0.2.1]`, true},
		{`"john.doe"@example.org`, true},
		{`"john doe"@example.org`, true},
		{`"john doe@example.org`, false},
		{`" john doe"@example.org`, true},
		{`""@example.org`, true},

		// in the next two, s/example/spodhuis/ to get a real address, by explicit configuration not catchall
		{"\"a~`*&^$#_-={}'?b\"@example.org", true},
		{`"X'); DROP TABLE domains; DROP TABLE passwords; --"@example.org`, true},
	})
}

func TestEmailAddressOrUnqualified(t *testing.T) {
	iterateBoolPatternMatch(t, EmailAddressOrUnqualified, "EmailAddressOrUnqualified", []boolPatternMatch{
		{`john`, true},
		{`john@example.org`, true},
		{`john:`, false},
		{`"john:"`, true},
		{`"john:"@example.org`, true},
		{`#`, true},  // beware using for a comment
		{`;`, false}, // better comment character
		{`# foo`, false},
		{`"# foo"`, true},
	})
}
