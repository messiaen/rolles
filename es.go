package rolles

import (
	"fmt"

	"github.com/olivere/elastic"
	"github.com/spf13/pflag"
)

type ElasticsearchOptions struct {
	Url string
}

func NewDefaultElasticsearchOptions() ElasticsearchOptions {
	return ElasticsearchOptions{
		Url: "http://localhost:9200",
	}
}

func (o *ElasticsearchOptions) BindFlags(flags *pflag.FlagSet) {
	f := "es"
	flags.StringVar(&o.Url, f, o.Url, "Elasticsearch address")
}

func (o ElasticsearchOptions) NewClient() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(o.Url), elastic.SetSniff(false))
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to Elasticsearch at %s -- %v", o.Url, err)
	}
	return client, nil
}
