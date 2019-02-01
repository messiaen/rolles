package rolles

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/messiaen/rolles/cmd"
	"github.com/olivere/elastic"
)

type TemplateOptions struct {
	EsOptions ElasticsearchOptions
	Dir       string
	Prefix    string
	Name      string
}

func NewDefaultTemplateOptions() TemplateOptions {
	return TemplateOptions{
		EsOptions: NewDefaultElasticsearchOptions(),
		Dir:       "./templates",
		Prefix:    "default",
		Name:      "",
	}
}

func toTempName(prefix, dir, fn string) string {
	name := fmt.Sprintf(strings.TrimPrefix(fn, fmt.Sprintf("%s/", dir)))
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_")
	return strings.Join([]string{prefix, strings.TrimSuffix(replacer.Replace(name), ".json")}, ".")
}

func (o TemplateOptions) GatherTemplates() (map[string]string, error) {
	if !cmd.FileExists(o.Dir) {
		return nil, fmt.Errorf("Template directory '%s' must exist", o.Dir)
	}
	dir, err := filepath.Abs(o.Dir)
	if err != nil {
		return nil, fmt.Errorf("Could not determine location of tempDir: %v", err)
	}

	var pairs map[string]string
	if o.Name != "" {
		pairs, err = gatherTemplatePattern(dir, o.Prefix, o.Name)
	} else {
		pairs, err = gatherAllTemplates(dir, o.Prefix)
	}
	if err != nil {
		return nil, err
	}
	return pairs, nil
}

func gatherAllTemplates(dir, prefix string) (map[string]string, error) {
	pairs := make(map[string]string, 0)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if dir == path {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		pairs[toTempName(prefix, dir, path)] = path
		return nil
	})
	return pairs, nil
}

func gatherTemplatePattern(dir, prefix, pattern string) (map[string]string, error) {
	matches, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		return nil, err
	}
	if len(matches) < 1 {
		matches, err = filepath.Glob(filepath.Join(dir, fmt.Sprintf("%s.json", pattern)))
	}
	if err != nil {
		return nil, err
	}

	if len(matches) < 1 {
		return nil, fmt.Errorf("No such template '%s' rooted at '%s'", pattern, dir)
	}
	pairs := make(map[string]string, 0)
	for _, m := range matches {
		pairs[toTempName(prefix, dir, m)] = m
	}
	return pairs, nil
}

func ReadTemplate(fn string) (string, error) {
	tempBytes, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", fmt.Errorf("Unable to read template at '%s' -- %v", fn, err)
	}
	tempStr := string(tempBytes)
	return tempStr, nil
}

func (o TemplateOptions) PutTemplateFromFilename(es *elastic.Client, ctx context.Context, name, fn string) error {
	tempStr, err := ReadTemplate(fn)
	if err != nil {
		return err
	}

	return o.PutTemplate(es, ctx, name, tempStr)
}

func (o TemplateOptions) PutTemplate(esClient *elastic.Client, ctx context.Context, name, t string) error {
	exists, err := esClient.IndexTemplateExists(name).Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to PUT template '%s' -- %v", name, err)
	}
	if exists {
		fmt.Println(fmt.Errorf("template '%s' already exists", name))
	}
	temp := make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(t), &temp)
	setIndexPatterns(temp, o.Prefix)
	setAliases(temp, o.Prefix)
	_, err = esClient.IndexPutTemplate(name).BodyJson(temp).Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to PUT template '%s' -- %v", name, err)
	}
	return nil
}

func setAliases(temp map[string]interface{}, prefix string) {
	aliases, has := temp["aliases"]
	if !has {
		return
	}
	newAliases := make(map[string]interface{}, 0)
	switch a := aliases.(type) {
	case map[string]interface{}:
		for aliasName, v := range a {
			newName := aliasName
			if strings.HasPrefix(newName, "*.") {
				newName = strings.TrimPrefix(newName, "*.")
			}
			newAliases[strings.Join([]string{prefix, newName}, ".")] = v
		}
		temp["aliases"] = newAliases
	}
}

func setIndexPatterns(temp map[string]interface{}, prefix string) {
	indexPats, has := temp["index_patterns"]
	if !has {
		return
	}
	switch s := indexPats.(type) {
	case []interface{}:
		pats := make([]string, len(s))
		for i, p := range s {
			tempPat := p.(string)
			if strings.HasPrefix(tempPat, "*.") {
				tempPat = strings.TrimPrefix(tempPat, "*.")
			}
			pats[i] = strings.Join([]string{prefix, tempPat}, ".")
		}
		temp["index_patterns"] = pats
	}
}

func (o TemplateOptions) DelTemplate(esClient *elastic.Client, ctx context.Context, name string) error {
	exists, err := esClient.IndexTemplateExists(name).Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to DELETE template '%s' -- %v", name, err)
	}
	if !exists {
		return nil
	}
	_, err = esClient.IndexDeleteTemplate(name).Do(ctx)
	if err != nil {
		return fmt.Errorf("Failed to DELETE template '%s' -- %v", name, err)
	}
	return nil
}
