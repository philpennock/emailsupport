// © Phil Pennock 2014.  See LICENSE file for licensing.

// +build rfc2822

package emailsupport

import (
	"testing"
)

func TestEmailLHSRFC2822(t *testing.T) {
	iterateBoolPatternMatch(t, EmailLHS, "EmailLHS/RFC2822", []boolPatternMatch{
		{"\"\x01\x02\"", true},
		{"\"\x01\\\x02\"", true},
	})
}
