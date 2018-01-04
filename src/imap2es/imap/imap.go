/*
   Copyright (C) 2018 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       imap2es
   Description   An IMAP to Elasticsearch indexer
   License       GPL version 2 (see LICENSE for details)
*/

package imap

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/go-ini/ini"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Connect(section *ini.Section) (*client.Client, error) {
	var c *client.Client
	var err error

	host := section.Key("host").String()
	scheme := section.Key("scheme").String()
	if scheme == "imap" {
		port := section.Key("port").MustString("143")
		c, err = client.Dial(host + ":" + port)
	} else if scheme == "imaps" {
		port := section.Key("port").MustString("993")
		c, err = client.DialTLS(host+":"+port, nil)
	} else {
		err = &errorString{"Scheme not supported: " + scheme}
	}

	if err != nil {
		return nil, err
	}

	caps, err := c.Capability()
	if caps["STARTTLS"] {
		c.StartTLS(nil)
	} else if err != nil {
		return nil, err
	}

	user := section.Key("user").String()
	password := section.Key("password").String()

	err = c.Login(user, password)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func RetrieveFolders(c *client.Client, folder string) []string {
	var folders []string
	mailboxes := make(chan *imap.MailboxInfo)

	go func() {
		c.List("", folder, mailboxes)
	}()

	for mailbox := range mailboxes {
		folders = append(folders, mailbox.Name)
	}

	return folders
}

func RetrieveMessages(c *client.Client, folder string, start, end uint32) ([]*imap.Message, error) {
	var emails []*imap.Message
	_, err := c.Select(folder, true)
	if err != nil {
		return nil, err
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(start, end)

	messages := make(chan *imap.Message, (end - start + 1))
	err = c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	if err != nil {
		return nil, err
	}

	for msg := range messages {
		emails = append(emails, msg)
	}

	return emails, nil
}

func CountMessages(c *client.Client, folder string) (uint32, error) {
	mbox, err := c.Select(folder, true)
	if err != nil {
		return 0, err
	}

	count := mbox.Messages

	return count, nil
}