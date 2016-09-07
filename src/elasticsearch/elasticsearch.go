package elasticsearch

import (
	"fmt"
	"github.com/go-ini/ini"
    "net/url"
    //"gopkg.in/olivere/elastic.v2"
)

func Connect(section *ini.Section) {
	
    url := &url.URL{
        Host: section.Key("host").String() +  ":" + section.Key("port").MustString("9200"),
        Scheme: section.Key("scheme").String(),
    }

    fmt.Println(url)
}
