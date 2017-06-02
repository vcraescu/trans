package loaders

import (
	"testing"
	"fmt"
	"math/rand"
	"path/filepath"
	"os"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestYamlLoader_Load(t *testing.T) {
	yaml := "foo:\n  bar:\n    foo: This is a test\nbar:\n  foo: This is foo value"

	filename := tmpName("test", ".yml")
	defer os.Remove(filename)

	ioutil.WriteFile(filename, []byte(yaml), 0755)

	loader := YamlLoader{}
	bag, err := loader.Load(filename)
	assert.Nil(t, err)

	r, ok := bag.Get("foo.bar.foo")
	assert.True(t, ok)
	assert.Equal(t, "This is a test", r)
}

func tmpName(prefix, suffix string) string {
	return filepath.Join(os.TempDir(), prefix + fmt.Sprintf("%x", rand.Int63()) + suffix)
}
