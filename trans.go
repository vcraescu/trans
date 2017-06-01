package trans

import (
	"io/ioutil"
	"errors"
	"strings"
	"github.com/vcraescu/trans/loaders"
	"path"
	"log"
	"github.com/vcraescu/trans/store"
	"github.com/vcraescu/trans/template"
)

type Locale string

type Params map[string]string

type Translator struct {
	FallbackLocale Locale
	Locale Locale
	catalog map[Locale]store.Bag
}

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

func (t *Translator) Load(dirname string) error {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return err
	}

	catalog := make(map[Locale]store.Bag)

	for _, file := range files {
		domain, locale, extension, err := parseTransFilename(file.Name())
		if err != nil {
			log.Println(err)
			continue
		}

		loader, err := loaders.NewByExtension(extension)
		if err != nil {
			log.Println(err)
			continue
		}

		bag, err := loader.Load(path.Join(dirname, file.Name()))
		if err != nil {
			log.Println(err)
			continue
		}

		domainBag := store.NewDataBag()
		domainBag.Set(domain, bag.All())
		catalog[locale] = domainBag
	}

	t.catalog = catalog

	return nil
}

func parseTransFilename(filename string) (domain string, locale Locale, extension string, err error) {
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
		locale, ok := option.(Locale)
		if ok {
			return locale, ok
		}
	}

	return Locale(""), true
}

func extractParams(options ...interface{}) (Params, bool) {
	for _, option := range options {
		params, ok := option.(Params)
		if ok {
			return params, ok
		}
	}

	return nil, false
}
