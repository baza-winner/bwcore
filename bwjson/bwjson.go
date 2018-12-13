// Package bwjson содержит bw-дополнение для package encoding/json
package bwjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
