package rolles

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/olivere/elastic"
)

type AliasOptions struct {
	EsOptions ElasticsearchOptions
	Config    string
	Prefix    string
	Name      string
}

type AliasConfiguration struct {
	BaseName           string             `json:"base_name"`
	RefDateFormat      string             `json:"ref_date_format"`
	RolloverConditions RolloverConditions `json:"rollover_conditions,omitempty"`
}

type RolloverConditions struct {
	MaxAge  string `json:"max_age"`
	MaxDocs int64  `json:"max_docs"`
}

func NewDefaultAliasOptions() AliasOptions {
	return AliasOptions{
		EsOptions: NewDefaultElasticsearchOptions(),
		Config:    "./index_conf.json",
		Prefix:    "default",
		Name:      "",
	}
}

func (o AliasOptions) GatherAliasCfgsFromFile() ([]AliasConfiguration, error) {
	cfgs, err := o.ReadAliasConfigs()
	if err != nil {
		return nil, err
	}
	return o.GatherAliasConfigs(cfgs), nil
}

func (o AliasOptions) GatherAliasConfigs(cfgs []AliasConfiguration) []AliasConfiguration {
	if o.Name == "" {
		return cfgs
	}
	gathered := make([]AliasConfiguration, 0)
	for _, c := range cfgs {
		// TODO add pattern matching
		if c.BaseName == o.Name {
			gathered = append(gathered, c)
		}
	}
	return gathered
}

func (o AliasOptions) ReadAliasConfigs() ([]AliasConfiguration, error) {
	cfgBytes, err := ioutil.ReadFile(o.Config)
	if err != nil {
		return nil, fmt.Errorf("Error reading alias config file '%s' -- %v", o.Config, err)
	}

	var cfg = []AliasConfiguration{}
	err = json.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse alias config '%s' -- %v", o.Config, err)
	}
	return cfg, nil
}

func (o AliasOptions) WriteAliasName(c AliasConfiguration) string {
	// TODO the '-w-' in the Write alias name should be configurable
	return fmt.Sprintf("%s.%s-w", o.Prefix, c.BaseName)
}

func (o AliasOptions) IndexName(c AliasConfiguration, time time.Time) string {
	timeStr := time.Format(c.RefDateFormat)
	return fmt.Sprintf("%s.%s-w-%s-00001", o.Prefix, c.BaseName, timeStr)
}

func (o AliasOptions) CollectAliases(es *elastic.Client, ctx context.Context) (map[string][]string, error) {
	rows, err := es.CatAliases().Do(ctx)
	if err != nil {
		return nil, err
	}
	aliases := make(map[string][]string)
	for _, r := range rows {
		if _, ok := aliases[r.Alias]; !ok {
			aliases[r.Alias] = make([]string, 0)
		}
		aliases[r.Alias] = append(aliases[r.Alias], r.Index)
	}
	return aliases, nil
}

func (o AliasOptions) DelAlias(esClient *elastic.Client, ctx context.Context, c AliasConfiguration, t time.Time) error {
	aliases, err := o.CollectAliases(esClient, ctx)
	if err != nil {
		return err
	}

	if indices, has := aliases[o.WriteAliasName(c)]; has && len(indices) == 1 {
		_, err = esClient.DeleteIndex(indices[0]).Do(ctx)
		if err != nil {
			return err
		}
	} else if !has {
		return nil
	} else {
		return fmt.Errorf("Aliases appear to be misconfigured")
	}

	return nil
}

func (o AliasOptions) PutAlias(esClient *elastic.Client, ctx context.Context, c AliasConfiguration, t time.Time) error {
	_, err := esClient.Aliases().Alias(o.WriteAliasName(c)).Do(ctx)
	if err != nil && !elastic.IsNotFound(err) {
		return err
	}
	_, err = esClient.CreateIndex(o.IndexName(c, t)).Do(ctx)
	if err != nil && !elastic.IsStatusCode(err, 400) {
		return err
	}
	_, err = esClient.Alias().Add(o.IndexName(c, t), o.WriteAliasName(c)).Do(ctx)
	if err != nil && !elastic.IsStatusCode(err, 400) {
		return err
	}
	return nil
}

func (o AliasOptions) Rollover(esClient *elastic.Client, ctx context.Context, c AliasConfiguration) error {
	_, err := esClient.Aliases().Alias(o.WriteAliasName(c)).Do(ctx)
	if err != nil && !elastic.IsNotFound(err) {
		return err
	}

	rollSvc := esClient.RolloverIndex(o.WriteAliasName(c)).
		AddMaxIndexAgeCondition(c.RolloverConditions.MaxAge).
		AddMaxIndexDocsCondition(c.RolloverConditions.MaxDocs)

	_, dryErr := rollSvc.DryRun(true).Do(ctx)

	if dryErr != nil {
		return fmt.Errorf("DryRun rollover of %s failed -- %v", o.WriteAliasName(c), dryErr)
	}

	rollRes, rollErr := rollSvc.DryRun(false).Do(ctx)
	if rollErr != nil {
		return fmt.Errorf("Rollover of %s failed -- %v\n", o.WriteAliasName(c), dryErr)
	} else if !rollRes.Acknowledged || !rollRes.ShardsAcknowledged || !rollRes.RolledOver {
		return fmt.Errorf("Rollover of %s failed -- %v\n", o.WriteAliasName(c), rollRes)
	}
	return nil
}
