package loaders

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewByExtension(t *testing.T) {
	RegisterLoader("yml", &YamlLoader{})

	{
		_, err := NewByExtension("yml")
		assert.Nil(t, err)
	}

	{
		_, err := NewByExtension("csv")
		assert.NotNil(t, err)
	}
}
