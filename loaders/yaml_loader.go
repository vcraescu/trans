package loaders

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/vcraescu/databag"
)

// YAML translation loader
type YamlLoader struct {
}

// Loads the translation keys from YML file
func (l *YamlLoader) Load(filename string) (databag.Bag, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make(map[interface{}]interface{})

	err = yaml.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}

	return databag.NewDataBagFrom(data), nil
}
