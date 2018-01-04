/*
   Copyright (C) 2018 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       imap2es
   Description   An IMAP to Elasticsearch indexer
   License       GPL version 2 (see LICENSE for details)
*/

package main

import (
	"imap2es/elasticsearch"
	"fmt"
	"github.com/gaal/go-options/options"
	"github.com/go-ini/ini"
	"imap2es/imap"
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

	imapcfg, _ := cfg.GetSection("imap")
	imapclient, err := imap.Connect(imapcfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer imapclient.Logout()

	folders := imap.RetrieveFolders(imapclient, "*")
	for _, mailbox := range folders {
		count, err := imap.CountMessages(imapclient, mailbox)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(mailbox+ ":", count)

		messages, err := imap.RetrieveMessages(imapclient, mailbox, 1, count)
		for _, message := range messages {
			fmt.Println("|--", message.Envelope.Subject)

			err = elasticsearch.Index(esclient, escfg.Key("index").String(), message)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}
