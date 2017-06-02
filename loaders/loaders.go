package loaders

import (
	"errors"
	"fmt"
	"github.com/vcraescu/databag"
)

type Loader interface {
	Load(filename string) (databag.Bag, error)
}

var registeredLoaders = make(map[string]Loader)

// Register a new translation loader based on file extension
func RegisterLoader(extension string, loader Loader) {
	registeredLoaders[extension] = loader
}

// Create a new translation loader based on the given file translation
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
