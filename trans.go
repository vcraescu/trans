// Simple translator similar to the one from Symfony2/3 framework
package trans

import (
	"io/ioutil"
	"errors"
	"strings"
	"github.com/vcraescu/trans/loaders"
	"path"
	"log"
	"github.com/vcraescu/trans/template"
	"github.com/vcraescu/databag"
)

// Translation locale.
type Locale string

// Translation placeholders.
type Params map[string]string

type Translator struct {
	// if there is no translation key for Locale, it will fallback to this locale
	FallbackLocale Locale
	Locale Locale
	catalog map[Locale]databag.Bag
}

// Translate a given key based on the given options.
// If the translation does not exists, the given key will be returned instead.
func (t *Translator) T(key string, options ...interface{}) string {
	locale := t.Locale

	optLocale, ok := extractLocale(options...)
	if ok {
		locale = optLocale
	}

	lc, ok := t.catalog[locale]
	if !ok && locale != t.FallbackLocale {
		lc, ok = t.catalog[t.FallbackLocale]
	}

	if !ok {
		return key
	}

	v, ok := lc.Get(key)
	if !ok {
		lc, ok = t.catalog[t.FallbackLocale]
	}

	if !ok {
		return key
	}

	v, ok = lc.Get(key)
	if !ok {
		return key
	}

	params, ok := extractParams(options...)
	str, _ := v.(string)

	return template.Parse(str, params)
}

// Loads all the translation files from the given directory.
//
// Translation files must have this format:
//  domain.locale.ext
// e.g:
//
// messages.en.yml
//
// validations.en.yml
//
// If there is no registered loader for the extension, the file will be ignored.
func (t *Translator) Load(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	if t.catalog == nil {
		t.catalog = make(map[Locale]databag.Bag)
	}

	for _, file := range files {
		bag, locale, err := loadCatalog(path.Join(dirname, file.Name()))
		if err != nil {
			log.Println(err)
			continue
		}

		if _, ok := t.catalog[locale]; !ok {
			t.catalog[locale] = bag
			continue
		}

		t.catalog[locale].Merge(bag)
	}

	return nil
}

func loadCatalog(path string) (databag.Bag, Locale, error) {
	domain, locale, extension, err := parseTransFilename(path)
	if err != nil {
		return nil, "", err
	}

	loader, err := loaders.NewByExtension(extension)
	if err != nil {
		return nil, "", err
	}

	bag, err := loader.Load(path)
	if err != nil {
		return nil, "", err
	}

	domainBag := databag.NewDataBag()
	domainBag.Set(domain, bag.All())

	return domainBag, locale, nil
}

func parseTransFilename(filename string) (domain string, locale Locale, extension string, err error) {
	filename = path.Base(filename)
	parts := strings.Split(filename, ".")
	if len(parts) < 3 {
		return "", "", "", errors.New("Invalid filename")
	}

	l := parts[len(parts)-2:]
	domain = strings.Join(parts[0:len(parts)-2], ".")
	return domain, Locale(l[0]), l[1], nil
}

func extractLocale(options ...interface{}) (Locale, bool) {
	for _, option := range options {
		if locale, ok := option.(Locale); ok {
			return locale, ok
		}
	}

	return Locale(""), false
}

func extractParams(options ...interface{}) (Params, bool) {
	for _, option := range options {
		if params, ok := option.(Params); ok {
			return params, ok
		}
	}

	return nil, false
}
