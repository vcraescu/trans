package store

import (
	"strings"
	"reflect"
)

const defaultNamespaceSep = "."

type Bag interface {
	Get(name string) (interface{}, bool)
	Set(name string, value interface{})
	All() map[interface{}]interface{}
}

type DataBag struct {
	data map[interface{}]interface{}
	namespaceSep string
}

func NewDataBag() *DataBag {
	return NewDataBagFrom(make(map[interface{}]interface{}))
}

func NewDataBagFrom(data map[interface{}]interface{}) *DataBag {
	return &DataBag{
		data: data,
		namespaceSep: defaultNamespaceSep,
	}
}

func namespaceSplitter(name string, sep string) []string {
	return strings.Split(name, sep)
}

func (d *DataBag) Get(name string) (interface{}, bool) {
	keys := namespaceSplitter(name, d.namespaceSep)

	count := len(keys)
	if count == 1 {
		v, ok := d.data[name]
		return v, ok
	}

	cp := d.data
	for i, key := range keys {
		v, ok := cp[key]
		if !ok || i == count-1 {
			return v, ok
		}

		_, ok = cp[key].(map[interface{}]interface{})
		if !ok {
			return nil, ok
		}

		cp = cp[key].(map[interface{}]interface{})
	}

	return nil, false
}

func (d *DataBag) Set(name string, value interface{}) {
	keys := namespaceSplitter(name, d.namespaceSep)

	count := len(keys)
	if count == 1 {
		d.data[name] = value
		return
	}

	cp := d.data
	for i := 0; i < count; i++ {
		key := keys[i]
		if i == count-1 {
			cp[key] = value
			break
		}

		_, ok := cp[key]
		if !ok || reflect.TypeOf(cp[key]).Kind() != reflect.Map {
			cp[key] = make(map[interface{}]interface{})
		}

		cp = cp[key].(map[interface{}]interface{})
	}
}

func (d *DataBag) All() map[interface{}]interface{} {
	return d.data
}
