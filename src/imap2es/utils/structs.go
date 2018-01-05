/*
   Copyright (C) 2018 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       imap2es
   Description   An IMAP to Elasticsearch indexer
   License       GPL version 2 (see LICENSE for details)
*/

package utils

import "github.com/emersion/go-imap"

type Message struct {
	From      []*imap.Address
	To        []*imap.Address
	Subject   string
	Body      string
	MessageId string
}
