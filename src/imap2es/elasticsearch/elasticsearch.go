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
   "settings":{
      "index":{
         "number_of_shards":5,
         "number_of_replicas":1
      }
   },
   "mappings":{
      "messages":{
         "properties":{
            "From":{
               "properties":{
                  "Address":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  },
                  "Name":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  }
               }
            },
            "To":{
               "properties":{
                  "Address":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  },
                  "Name":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  }
               }
            },
            "CC":{
               "properties":{
                  "Address":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  },
                  "Name":{
                     "type":"text",
                     "fields":{
                        "keyword":{
                           "ignore_above":256,
                           "type":"keyword"
                        }
                     }
                  }
               }
            },
            "Folder":{
               "type":"text",
               "fields":{
                  "keyword":{
                     "ignore_above":256,
                     "type":"keyword"
                  }
               }
            },
            "Body":{
               "type":"text",
               "fields":{
                  "keyword":{
                     "ignore_above":256,
                     "type":"keyword"
                  }
               }
            },
            "Date":{
               "type":"date"
            },
            "Subject":{
               "type":"text",
               "fields":{
                  "keyword":{
                     "ignore_above":256,
                     "type":"keyword"
                  }
               }
            },
            "MessageId":{
               "type":"text",
               "fields":{
                  "keyword":{
                     "ignore_above":256,
                     "type":"keyword"
                  }
               }
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

	uri := &url.URL{
		Host:   section.Key("host").String() + ":" + section.Key("port").MustString("9200"),
		Scheme: section.Key("scheme").String(),
	}

	// Create a client
	client, err := elastic.NewClient(elastic.SetURL(uri.String()))
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
	var err error

	err = createIndexIfNotExists(client, index)
	if err != nil {
		return err
	}

	_, err = client.Index().Index(index).Type("messages").BodyJson(message).Do(context.Background())

	if err != nil {
		return err
	} else {
		return nil
	}
}
