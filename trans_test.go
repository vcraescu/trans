package trans

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"math/rand"
	"path"
	"strings"
)

func TestExtractParams(t *testing.T) {
	{
		params := Params{"foo": "bar"}
		actual, ok := extractParams([]string{"test"}, Locale("en"), params)
		assert.True(t, ok)
		assert.Exactly(t, params, actual)
	}
	{
		params1 := Params{"foo": "bar"}
		params2 := Params{"bar": "foo"}
		actual, ok := extractParams([]string{"test"}, Locale("en"), params1, params2)
		assert.True(t, ok)
		assert.Exactly(t, params1, actual)
	}
	{
		var params Params
		actual, ok := extractParams([]string{"test"}, Locale("en"))
		assert.False(t, ok)
		assert.Exactly(t, params, actual)
	}
}

func TestExtractLocale(t *testing.T) {
	{
		locale := Locale("en")
		actual, ok := extractLocale([]string{"test"}, Locale("en"), Params{})
		assert.True(t, ok)
		assert.Exactly(t, locale, actual)
	}
	{
		var locale Locale
		actual, ok := extractLocale([]string{"test"}, Params{})
		assert.False(t, ok)
		assert.Exactly(t, locale, actual)
	}
}

func TestParseTransFilename(t *testing.T)  {
	{
		domain, locale, extension, err := parseTransFilename("/foo.en/bar.yml/message.en.yml")
		assert.Nil(t, err)
		assert.Exactly(t, "message", domain)
		assert.Exactly(t, Locale("en"), locale)
		assert.Exactly(t, "yml", extension)
	}
	{
		_, _, _, err := parseTransFilename("message.yml")
		assert.NotNil(t, err)
	}
}

func TestLoadCatalog(t *testing.T) {
	yaml := "foo:\n  bar:\n    foo: This is a test\nbar:\n  foo: This is foo value"

	filename := tmpName("test", ".en.yml")
	defer os.Remove(filename)

	ioutil.WriteFile(filename, []byte(yaml), 0755)

	parts := strings.Split(path.Base(filename), ".")
	domain := parts[0]

	bag, locale, err := loadCatalog(filename)
	assert.Nil(t, err)
	assert.Exactly(t, Locale("en"), locale)

	actual, ok := bag.Get(domain + ".foo.bar.foo")
	assert.True(t, ok)
	assert.Exactly(t, actual, "This is a test")
}

func TestTranslatorLoadMerge(t *testing.T) {

	dirname1 := path.Join(os.TempDir(), "trans1")
	os.Mkdir(dirname1, 0755)
	defer os.RemoveAll(dirname1)

	dirname2 := path.Join(os.TempDir(), "trans2")
	os.Mkdir(dirname2, 0755)
	defer os.RemoveAll(dirname2)

	filename1 := path.Join(os.TempDir(), "/trans1/", "messages.en.yml")
	{
		yaml := "foo:\n  bar:\n    foo: This is a test\nbar:\n  foo: This is foo value"
		err := ioutil.WriteFile(filename1, []byte(yaml), 0755)
		assert.Nil(t, err)
	}

	filename2 := path.Join(os.TempDir(), "/trans2/", "messages.en.yml")
	{
		yaml := "bar:\n  foo: This is another foo value"
		err := ioutil.WriteFile(filename2, []byte(yaml), 0755)
		assert.Nil(t, err)
	}

	trans := Translator{Locale:Locale("en")}
	{
		err := trans.Load(dirname1)
		assert.Nil(t, err)

		err = trans.Load(dirname2)
		assert.Nil(t, err)

		{
			actual := trans.T("messages.foo.bar.foo")
			assert.Exactly(t, "This is a test", actual)
		}
		{
			actual := trans.T("messages.this.key.does.not.exists")
			assert.Exactly(t, "messages.this.key.does.not.exists", actual)
		}
		{
			actual := trans.T("messages.this.key.does.not.exists")
			assert.Exactly(t, "messages.this.key.does.not.exists", actual)
		}
		{
			actual := trans.T("messages.bar.foo")
			assert.Exactly(t, actual, "This is another foo value")
		}
	}
}

func tmpName(prefix, suffix string) string {
	return filepath.Join(os.TempDir(), prefix + fmt.Sprintf("%x", rand.Int63()) + suffix)
}

func ExampleTranslator_Load() {
	// you can load multiple translation files from different locations and all the translation keys will be merged
	// into the translator

	t := Translator{}
	t.Load("/path/to/translations")
	t.Load("/another/path/to/translations")
}

func ExampleTranslator_T() {
	// given we have loaded the following translation keys:
	// messages.core.flash.error = %product% product cannot be deleted
	// messages.core.buttons.ok_label = OK

	t := Translator{}
	t.Load("/path/to/translations")
	// and the following files exists into /path/to/translations folder
	// messages.en.yml
	// messages.de.yml

	text := t.T("messages.core.buttons.ok_label")

	fmt.Println(text)
	// => OK

	text = t.T("messages.core.flash.error", Params{"%product%": "T-Shirt"}, Locale("de"))

	fmt.Println(text)
	// => T-Shirt product cannot be deleted
}

