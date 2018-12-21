package bw

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/davecgh/go-spew/spew"
)

//go:generate stringer -type=ValPathItemType

// ============================================================================

var Spew spew.ConfigState

func init() {
	Spew = spew.ConfigState{SortKeys: true}
}

// ============================================================================

func Args(args ...interface{}) []interface{} {
	return args
}

// ============================================================================

type I interface {
	FmtString() string
	FmtArgs() []interface{}
}

// ============================================================================

type A struct {
	Fmt  string
	Args []interface{}
}

func (v A) FmtString() string {
	return v.Fmt
}

func (v A) FmtArgs() []interface{} {
	return v.Args
}

// ============================================================================

func Fmt(fmtString string, fmtArg ...interface{}) A {
	return A{fmtString, fmtArg}
}

// ============================================================================

func PluralWord(count int, word string, word1 string, word2_4 string, _word5more ...string) (result string) {
	var word5more string
	if _word5more != nil {
		word5more = _word5more[0]
	}
	if len(word5more) == 0 {
		word5more = word2_4
	}
	result = word5more
	decimal := count / 10 % 10
	if decimal != 1 {
		unit := count % 10
		if unit == 1 {
			result = word1
		} else if 2 <= unit && unit <= 4 {
			result = word2_4
		}
	}
	return word + result
}

// ============================================================================

func NormalIdx(idx int, len int) (result int, ok bool) {
	min := -len
	max := len - 1
	if min <= idx && idx <= max {
		if idx < 0 {
			idx = len + idx
		}
		result = idx
		ok = true
	} else {
		ok = false
	}
	return
}

// ============================================================================

// https://stackoverflow.com/questions/6878590/the-maximum-value-for-an-int-type-in-go

const (
	MaxUint8 = ^uint8(0)
	MinUint8 = 0
	MaxInt8  = int8(MaxUint8 >> 1)
	MinInt8  = -MaxInt8 - 1

	MaxUint16 = ^uint16(0)
	MinUint16 = 0
	MaxInt16  = int16(MaxUint16 >> 1)
	MinInt16  = -MaxInt16 - 1

	MaxUint32 = ^uint32(0)
	MinUint32 = 0
	MaxInt32  = int32(MaxUint32 >> 1)
	MinInt32  = -MaxInt32 - 1

	MaxUint64 = ^uint64(0)
	MinUint64 = 0
	MaxInt64  = int64(MaxUint64 >> 1)
	MinInt64  = -MaxInt64 - 1

	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

// ============================================================================

type ValPathItemType uint8

const (
	ValPathItemDef ValPathItemType = iota
	ValPathItemFileSpec
	ValPathItemIdx
	ValPathItemKey
	ValPathItemVar
	ValPathItemFunc
	ValPathItemPath
)

func (v ValPathItemType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

type ValPathItem struct {
	Type       ValPathItemType
	Idx        int
	Name       string
	Path       ValPath
	IsOptional bool
	DefaultVal interface{}
	FuncArg    interface{}
	Def        interface{}
}

func (vpi ValPathItem) keyForPathS() (result string) {
	result = vpi.Name
	runes := []rune(result)
	var needQuote bool
	if needQuote = len(runes) == 0; !needQuote {
		for i, r := range runes {
			if r == '_' || unicode.IsLetter(r) || i > 0 && ('0' <= r && r <= '9') {
				continue
			}
			needQuote = true
			break
		}
	}
	if needQuote {
		result = fmt.Sprintf("%q", result)
	}
	return
}

type ValPath []ValPathItem

func (v ValPath) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, i := range v {
		result = append(result, i)
	}
	return json.Marshal(result)
}

func (v ValPath) String() (result string) {
	ss := []string{}
	if len(v) == 0 {
		result = "."
	} else {
		for _, vpi := range v {
			var s string
			switch vpi.Type {
			case ValPathItemPath:
				s = "(" + vpi.Path.String() + ")"
			case ValPathItemKey:
				s = vpi.keyForPathS()
			case ValPathItemVar:
				s = "$" + vpi.keyForPathS()
			case ValPathItemIdx:
				s = strconv.FormatInt(int64(vpi.Idx), 10)
			// case ValPathItemHash:
			// 	s = "#"
			default:
				panic(Spew.Sprintf("%#v", vpi.Type))
			}
			if vpi.IsOptional {
				s += "?"
			}
			ss = append(ss, s)
		}
		result = strings.Join(ss, ".")
	}
	return
}

func (v ValPath) Clone() ValPath {
	// https://github.com/go101/go101/wiki/How-to-efficiently-clone-a-slice%3F
	return append(v[:0:0], v...)
}

func (v ValPath) AppendIdx(idx int) ValPath {
	// path := append(v[:0:0], v...) // https://github.com/go101/go101/wiki/How-to-efficiently-clone-a-slice%3F
	return append(v.Clone(), ValPathItem{Type: ValPathItemIdx, Idx: idx})
}

func (v ValPath) AppendKey(key string) ValPath {
	// path := append(v[:0:0], v...)
	return append(v.Clone(), ValPathItem{Type: ValPathItemKey, Name: key})
}

func (v ValPath) AppendVar(name string) ValPath {
	return append(v.Clone(), ValPathItem{Type: ValPathItemVar, Name: name})
}

// func (v ValPath) AppendHash(name string) ValPath {
// 	return append(v.Clone(), ValPathItem{Type: ValPathItemHash})
// }

func (v ValPath) Append(a ValPath) ValPath {
	return append(v.Clone(), a...)
}

// ============================================================================

type ValPathProvider interface {
	Path() (ValPath, error)
}

func (p ValPath) Path() (ValPath, error) { return p, nil }

// ============================================================================

type Val interface {
	PathVal(path ValPath, optVars ...map[string]interface{}) (result interface{}, err error)
	SetPathVal(val interface{}, path ValPath, optVars ...map[string]interface{}) (err error)
	MarshalJSON() ([]byte, error)
}

// ============================================================================

var braces = map[rune]rune{
	'(': ')',
	'{': '}',
	'<': '>',
	'[': ']',
}

func Braces() map[rune]rune {
	return braces
}

// ============================================================================

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsLetter(r rune) bool {
	return r == '_' || unicode.IsLetter(r)
}

func IsPunctOrSymbol(r rune) bool {
	return unicode.IsPunct(r) || unicode.IsSymbol(r)
}

// ============================================================================
