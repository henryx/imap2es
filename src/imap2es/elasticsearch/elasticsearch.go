/*
   Copyright (C) 2018 Enrico Bianchi (enrico.bianchi@gmail.com)
   Project       imap2es
   Description   An IMAP to Elasticsearch indexer
   License       GPL version 2 (see LICENSE for details)
*/

package elasticsearch

import (
	"context"
	"github.com/go-ini/ini"
	"gopkg.in/olivere/elastic.v6"
	"imap2es/utils"
	"net/url"
)

const mapping = `
{
    "settings": {
        "index": {
            "number_of_shards": 5,
            "number_of_replicas": 1
        }
    },
    "mappings": {
        "messages": {
            "properties": {
                "from": {
                    "type": "keyword"
                },
                "to": {
                    "type": "keyword"
                },
                "messageid": {
                    "type": "text"
                },
                "folder": {
                    "type": "text"
                },
                "date": {
                    "type": "date",
                    "format": "EEE, d MMM yyyy HH:mm:ss Z"
                }
            }
        }
    }
}
`

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func Connect(section *ini.Section) (*elastic.Client, error) {

	url := &url.URL{
		Host:   section.Key("host").String() + ":" + section.Key("port").MustString("9200"),
		Scheme: section.Key("scheme").String(),
	}

	// Create a client
	client, err := elastic.NewClient(elastic.SetURL(url.String()))
	if err != nil {
		return nil, err
	}

	return client, nil
}

func createIndexIfNotExists(client *elastic.Client, index string) error {
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		return err
	}

	if !exists {
		idx, err := client.CreateIndex(index).BodyString(mapping).Do(context.Background())
		if err != nil {
			// Handle error
			panic(err)
		}
		if !idx.Acknowledged {
			return &errorString{"Index not acknowledged"}
		}
	}

	return nil
}

func Index(client *elastic.Client, index string, message utils.Message) error {
	err := createIndexIfNotExists(client, index)
	if err != nil {
		return err
	}

	client.Index().Index(index).BodyJson(message).Do(context.Background())

	return nil
}
