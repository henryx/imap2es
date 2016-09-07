package main

import (
	"elasticsearch"
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
	var cfg *ini.File
	var err error

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
		cfg, err = ini.Load([]byte{}, opts.Get("cfg"))
	}

	if err != nil {
		fmt.Println("Error about reading config file:", err)
		os.Exit(1)
	}

	escfg, _ := cfg.GetSection("elasticsearch")
	esclient, err := elasticsearch.Connect(escfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = elasticsearch.Index(esclient, escfg.Key("index").String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
