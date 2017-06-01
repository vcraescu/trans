package loaders

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"errors"
	"fmt"
	"github.com/vcraescu/trans/store"
)

type Loader interface {
	Load(filename string) (store.Bag, error)
}

type YamlLoader struct {
}

func (l *YamlLoader) Load(filename string) (store.Bag, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make(map[interface{}]interface{})

	err = yaml.Unmarshal(contents, &data)
	if err != nil {
		return nil, err
	}

	return store.NewDataBagFrom(data), nil
}

var registeredLoaders = make(map[string]Loader)

func RegisterLoader(extension string, loader Loader) {
	registeredLoaders[extension] = loader
}

func NewByExtension(extension string) (Loader, error) {
	loader, ok := registeredLoaders[extension]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unknown extension type: %s", extension))
	}

	return loader, nil
}

func init() {
	yamlLoader := &YamlLoader{}

	RegisterLoader("yaml", yamlLoader)
	RegisterLoader("yml", yamlLoader)
}
