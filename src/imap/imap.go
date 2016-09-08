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

func listFolders(client *imap.Client, mailbox string) chan *imap.MailboxInfo {
	var rsp *imap.Response
	var search string

	ch := make(chan *imap.MailboxInfo)

	cmd, _ := imap.Wait(client.List("", "INBOX"))
	delim := cmd.Data[0].MailboxInfo().Delim

	if mailbox != "INBOX" && mailbox != "" {
		search = mailbox + delim + "%"
	} else {
		search = "%"
	}

	cmd, _ = imap.Wait(client.List("", search))

	go func() {
		for _, rsp = range cmd.Data {
			ch <- rsp.MailboxInfo()
		}
		close(ch)
	}()

	return ch
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

func RetrieveFolders(client *imap.Client, folder string) []string {
	var folders []string

	for mailbox := range listFolders(client, folder) {
		folders = append(folders, mailbox.Name)

		if mailbox.Attrs["\\Haschildren"] == true {
			subfolders := RetrieveFolders(client, mailbox.Name)
			for _, subfolder := range subfolders {
				folders = append(folders, subfolder)
			}
		}
	}
	return folders
}

func MessageCount(client *imap.Client, folder string) int {
	client.Select(folder, true)

	return client.Mailbox.Messages
}
