// Package bwjson содержит bw-дополнение для package encoding/json
package bwjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwos"
)

// Pretty - wrapper для json.MarshalIndent
func Pretty(v interface{}) (result string) {
	if bytes, err := json.MarshalIndent(v, "", "  "); err != nil {
		panic(fmt.Sprintf("Pretty: failed to encode to json value %#v", v))
	} else {
		result = string(bytes[:])
	}
	return
}

func S(v interface{}) (result string) {
	if bytes, err := json.Marshal(v); err != nil {
		panic(fmt.Sprintf("S: failed to encode to json value %#v", v))
	} else {
		result = string(bytes[:])
	}
	return
}

func FromFile(fileSpec string) (val interface{}, err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(fileSpec); err != nil {
		return
	}
	defer jsonFile.Close()
	var bytes []byte
	if bytes, err = ioutil.ReadAll(jsonFile); err != nil {
		err = bwerr.Refine(err, "bwjson.<ansiFunc>FromFile<ansi>(<ansiPath>%s<ansi>): {Error}", bwos.ShortenFileSpec(fileSpec))
		return
	}
	if err = json.Unmarshal(bytes, &val); err != nil {
		err = bwerr.Refine(err, "bwjson.<ansiFunc>FromFile<ansi>(<ansiPath>%s<ansi>): {Error}", bwos.ShortenFileSpec(fileSpec))
		return
	}
	return
}

func ToFile(fileSpec string, val interface{}) (err error) {
	var bytes []byte
	if bytes, err = json.MarshalIndent(val, "", "  "); err != nil {
		return bwerr.Refine(err, "bwjson.ToFile(%s): {Error}", bwos.ShortenFileSpec(fileSpec))
	}
	err = ioutil.WriteFile(fileSpec, bytes, 0644)
	return
}

func MarshalJSON(val interface{}) ([]byte, error) {
	var value = reflect.ValueOf(val)
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return json.Marshal(val)
	} else {
		result := map[string]interface{}{}
	FIELD:
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			for fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					continue FIELD
				}
				fieldValue = fieldValue.Elem()
			}
			// kind := fieldValue.Kind()
			switch fieldValue.Kind() {
			case reflect.Slice, reflect.Map, reflect.String:
				if fieldValue.Len() == 0 {
					continue FIELD
				}
			case reflect.Interface:
				if fieldValue.IsNil() {
					continue FIELD
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if fieldValue.Int() == 0 {
					continue FIELD
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if fieldValue.Uint() == 0 {
					continue FIELD
				}
			case reflect.Float32, reflect.Float64:
				if fieldValue.Float() == 0 {
					continue FIELD
				}
			}
			// if (kind == reflect.Slice {
			// 	if fieldValue.Len() == 0 {
			// 		continue FIELD
			// 	}
			// } else if kind == reflect.Map {
			// 	if fieldValue.Len() == 0 {
			// 		continue FIELD
			// 	}
			// } else if kind == reflect.String {
			// 	if fieldValue.Len() == 0 {
			// 		continue FIELD
			// 	}
			// }
			result[value.Type().Field(i).Name] = fieldValue
		}
		return json.Marshal(result)
	}
}
