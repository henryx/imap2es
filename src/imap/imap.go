package imap

import (
	"github.com/go-ini/ini"
	"github.com/mxk/go-imap/imap"
)

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Connect(section *ini.Section) (*imap.Client, error) {
	var client *imap.Client
	var err error

	host := section.Key("host").String()
	scheme := section.Key("scheme").String()
	if scheme == "imap" {
		port := section.Key("port").MustString("143")
		client, err = imap.Dial(host + ":" + port)
	} else if scheme == "imaps" {
		port := section.Key("port").MustString("993")
		client, err = imap.DialTLS(host+":"+port, nil)
	} else {
		err = &errorString{"Scheme not supported: " + scheme}
	}

	if err != nil {
		return nil, err
	}

	if client.Caps["STARTTLS"] {
		client.StartTLS(nil)
	}

	if client.State() == imap.Login {
		user := section.Key("user").String()
		password := section.Key("password").String()

		_, err := client.Login(user, password)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func Mailboxes(client *imap.Client) chan *imap.MailboxInfo {
	var rsp *imap.Response
	ch := make(chan *imap.MailboxInfo)
	cmd, _ := imap.Wait(client.List("", "%"))

	go func() {
		for _, rsp = range cmd.Data {
			ch <- rsp.MailboxInfo()
		}
		close(ch)
	}()

	return ch
}
