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
	"github.com/emersion/go-message/mail"
	"github.com/go-ini/ini"
	"imap2es/utils"
	"io"
	"io/ioutil"
	"log"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func parseMessage(msg *imap.Message) utils.Message {
	retval := utils.Message{}

	section, _ := imap.ParseBodySectionName("BODY[]")
	raw := msg.GetBody(section)

	reader, err := mail.CreateReader(raw)
	if err != nil {
		log.Fatal(err)
	}

	header := reader.Header
	if date, err := header.Date(); err == nil {
		retval.Date = date
	}

	if from, err := header.AddressList("From"); err == nil {
		retval.From = from
	}

	if to, err := header.AddressList("To"); err == nil {
		retval.To = to
	}

	if subject, err := header.Subject(); err == nil {
		retval.Subject = subject
	}

	body := ""
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		bodypart, _ := ioutil.ReadAll(part.Body)
		body += string(bodypart)
	}

	retval.Body = body

	return retval
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

func RetrieveMessages(c *client.Client, folder string, start, end uint32) ([]utils.Message, error) {
	var emails []utils.Message
	_, err := c.Select(folder, true)

	if start > end {
		return emails, nil
	}

	if err != nil {
		return nil, err
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(start, end)

	messages := make(chan *imap.Message, (end - start + 1))
	err = c.Fetch(seqset, []imap.FetchItem{"BODY[]"}, messages)
	if err != nil {
		return nil, err
	}

	for msg := range messages {
		message := parseMessage(msg)
		emails = append(emails, message)
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
