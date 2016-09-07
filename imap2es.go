package main

import (
	"fmt"
	"github.com/gaal/go-options/options"
	"github.com/go-ini/ini"
	"os"
)

const VERSION = "0.0.0"

const SPEC = `
An IMAP to Elasticsearch indexer
Usage: imap2es [OPTIONS]
--
h,help                Print this help
v,version             Print version
c,cfg=                Set the configuration file
`

func main() {
	s := options.NewOptions(SPEC)

	// Check if options isn't passed
	if len(os.Args[1:]) <= 0 {
		s.PrintUsageAndExit("No option specified")
	}
	opts := s.Parse(os.Args[1:])

	// Print version and exit
	if opts.GetBool("version") {
		fmt.Println("imap2es " + VERSION)
		os.Exit(0)
	}

	// Read the configuration file
	if opts.GetBool("cfg") {
		cfg, err = ini.Load([]byte{}, opts.Get("backup"))
	}

	if err != nil {
		fmt.Println("Error about reading config file:", err)
		os.Exit(1)
	}
}