package elasticsearch

import (
	"github.com/go-ini/ini"
	"gopkg.in/olivere/elastic.v2"
	"net/url"
)

const mapping = `{
    "settings":{
        "number_of_shards":5,
        "number_of_replicas":1
    }
}`

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
	exists, err := client.IndexExists(index).Do()
	if err != nil {
		return err
	}

	if !exists {
		idx, err := client.CreateIndex(index).BodyString(mapping).Do()
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

func Index(client *elastic.Client, index string) error {
	err := createIndexIfNotExists(client, index)
	if err != nil {
		return err
	}
	return nil
}
