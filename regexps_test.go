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

func TestIpv4Octets(t *testing.T) {
	iterateBoolPatternMatch(t, Ipv4Octet, "Ipv4Octet", []boolPatternMatch{
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

func TestIpv4Addresses(t *testing.T) {
	iterateBoolPatternMatch(t, Ipv4Address, "Ipv4Address", []boolPatternMatch{
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

func TestIpv4Netblocks(t *testing.T) {
	iterateBoolPatternMatch(t, Ipv4Netblock, "Ipv4Netblock", []boolPatternMatch{
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
func TestIpv6AddressesFromEmitTester(t *testing.T) {
	iterateBoolPatternMatch(t, Ipv6Address, "Ipv4Address", []boolPatternMatch{
		{"::", true},
		{"::1", true},
		{"fe02::1", true},
		{"::ffff:192.0.2.1", true},
		{"2001:DB8::42", true},
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

func TestIpv6Netblocks(t *testing.T) {
	iterateBoolPatternMatch(t, Ipv6Netblock, "Ipv6Netblock", []boolPatternMatch{
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

func TestDebug(t *testing.T) {
	t.Logf("\nATEXT: %s\nQTEXT: %s\nLHS:\n%s\nDomain: %s\n", txtAText, txtQText, txtEmailLHS, txtEmailDomain)
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
		{`"X'); DROP TABLE domains; DROP TABLE passwords; --"`, true},
		{`"<script>alert('Boo!')</script>"`, true},
	})
}
