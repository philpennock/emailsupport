package main

import (
	"fmt"

	"github.com/philpennock/emailsupport"
)

// TODO:
// Flags:
//  • / surrounds and replace internal "/" with "\/" (as emitted)
//  • replace "`" with "\x60"
//  • anchor
//  • name the pattern extractors (lhs, domain, ipv4, ipv6, etc)
func main() {
	fmt.Println(emailsupport.TxtEmailAddress)
}
