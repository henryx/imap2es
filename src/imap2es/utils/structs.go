/*
   Copyright (C) 2018 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       imap2es
   Description   An IMAP to Elasticsearch indexer
   License       GPL version 2 (see LICENSE for details)
*/

package utils

import (
	"github.com/emersion/go-message/mail"
	"time"
)

type Message struct {
	From      []*mail.Address
	To        []*mail.Address
	Subject   string
	Date      time.Time
	Body      string
	MessageId string
}
