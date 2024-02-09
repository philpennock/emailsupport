package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/philpennock/emailsupport"
)

const (
	// sysexits.h
	EX_USAGE   = 64
	EX_DATAERR = 65
	EX_NOINPUT = 66
)

// checker is not safe to share between threads: it prints to stdout, and thus
// we don't bother making the count field safe either.
type checker struct {
	okay  bool
	count int
}

func NewChecker() *checker { return &checker{okay: true} }

func (c *checker) IsEmailAddress(text string) {
	if emailsupport.EmailAddress.MatchString(text) {
		fmt.Printf("OK: %q\n", text)
	} else {
		fmt.Printf("FAIL: %q\n", text)
		c.okay = false
	}
	c.count += 1
}

func (c *checker) Stream(in io.Reader) {
	rdr := bufio.NewScanner(in)
	for rdr.Scan() {
		line := rdr.Text()
		if len(line) == 0 {
			continue
		}
		c.IsEmailAddress(line)
	}
}

func main() {
	ourName := filepath.Base(os.Args[0])
	inputFile := flag.String("file", "", "read addresses from file, one per line (no comments, no exceptions except completely blank lines)")
	flag.Parse()

	check := NewChecker()

	if *inputFile != "" {
		if len(flag.Args()) > 0 {
			fmt.Fprintf(os.Stderr, "%s: can't take parameters if using -file\n", ourName)
			os.Exit(EX_USAGE)
		}
		if *inputFile == "-" {
			check.Stream(os.Stdin)
		} else {
			fh, err := os.Open(*inputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: opening %q failed: %v\n", ourName, *inputFile, err)
				os.Exit(EX_NOINPUT)
			}
			func(in *os.File) {
				defer in.Close()
				check.Stream(fh)
			}(fh)
		}
	} else {
		if len(flag.Args()) == 0 {
			fmt.Fprintf(os.Stderr, "%s: need at lease one parameter to check email addresses\n", ourName)
			os.Exit(EX_USAGE)
		}

		for _, param := range flag.Args() {
			check.IsEmailAddress(param)
		}
	}

	if !check.okay {
		os.Exit(1)
	}

	if check.count == 0 {
		// given an empty input file, or one containing only blank lines?
		fmt.Fprintf(os.Stderr, "%s: did not check any addresses\n", ourName)
		os.Exit(EX_DATAERR)
	}
}
