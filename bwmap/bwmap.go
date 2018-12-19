/*
Предоставялет функции для работы со структурой map[string]interface{}.
*/
package bwmap

import (
	"reflect"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwset"
)

type I interface {
	HasKey(key string) bool
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	Keys(optFilter ...KeysFilter) (result []string)
	DelKey(key string)
	Map() map[string]interface{}
}

type M map[string]interface{}

var (
	ansiMustBeMap       string
	ansiMustBeMapString string
)

func init() {
	ansiMustBeMap = ansi.String("<ansiVar>m<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>map")
	ansiMustBeMapString = ansi.String("<ansiVar>m<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>map[string]")
}

func UnexpectedKeys(m interface{}, expected ...interface{}) (result bwset.String, err error) {
	if expected == nil {
		return
	}
	if o, ok := m.(*Ordered); ok {
		m = o.Map()
	}
	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		err = bwerr.From(ansiMustBeMap, m)
		return
	}
	for _, vk := range v.MapKeys() {
		if vk.Kind() != reflect.String {
			err = bwerr.From(ansiMustBeMapString, m)
			return
		}
		break
	}
	expectedKeys := map[string]struct{}{}
	for _, item := range expected {
		if item == nil {
			continue
		}
		if s, ok := item.(string); ok {
			expectedKeys[s] = struct{}{}
		} else if ss, ok := item.([]string); ok {
			for _, k := range ss {
				expectedKeys[k] = struct{}{}
			}
		} else if reflect.TypeOf(item).Kind() == reflect.Map {
			v := reflect.ValueOf(item)
			for _, vk := range v.MapKeys() {
				k := vk.String()
				expectedKeys[k] = struct{}{}
			}
		} else {
			err = bwerr.From("<ansiVar>expected<ansi> (<ansiVal>%+v<ansi>) neither <ansiVal>string<ansi> nor <ansiVal>[]string<ansi> nor <ansiVal>map[string]", expected)
			return
		}
	}
	result = bwset.String{}
	for _, vk := range v.MapKeys() {
		k := vk.String()
		if _, ok := expectedKeys[k]; !ok {
			result.Add(k)
		}
	}
	if len(result) == 0 {
		result = nil
	}
	return
}

func MustUnexpectedKeys(m interface{}, expected ...interface{}) (result bwset.String) {
	var err error
	if result, err = UnexpectedKeys(m, expected...); err != nil {
		bwerr.PanicA(bwerr.E{Depth: 1, Error: err})
	}
	return
}

type KeysFilter func(key string) (ok bool)

func (m M) Keys(optFilter ...KeysFilter) (result []string) {
	var filter KeysFilter
	if len(optFilter) > 0 {
		filter = optFilter[0]
	}
	for key, _ := range m {
		if filter == nil || filter(key) {
			result = append(result, key)
		}
	}
	return
}

func (m M) HasKey(key string) bool {
	_, exists := m[key]
	return exists
}

func (m M) Get(key string) (interface{}, bool) {
	val, exists := m[key]
	return val, exists
}

func (m M) Set(key string, val interface{}) {
	m[key] = val
}

func (m M) DelKey(key string) {
	delete(m, key)
}

func (m M) Map() map[string]interface{} {
	return m
}

// func ForEach(m map[string]interface{}, body ForEachBody) (err error) {
// }

// func CropMap(m interface{}, crop ...interface{}) {
// 	if unexpectedKeys, err := UnexpectedKeys(m, crop...); err != nil {
// 		bwerr.PanicErr(err)
// 	} else if unexpectedKeys != nil {
// 		for k, _ := range unexpectedKeys {
// 			v := reflect.ValueOf(m)
// 			v.SetMapIndex(reflect.ValueOf(k), reflect.Value{})
// 		}
// 	}
// }
