package bwval

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwset"
)

// ============================================================================

func Int(val interface{}) (result int, ok bool) {
	var needRecall bool
	ok = true
	switch t := val.(type) {
	case int8:
		result = int(t)
	case int16:
		result = int(t)
	case int32:
		result = int(t)
	case int:
		result = t
	case uint8:
		result = int(t)
	case uint16:
		result = int(t)
	case uint64:
		if ok = t <= uint64(bw.MaxInt); ok {
			result = int(t)
		}
	case uint:
		if ok = t <= uint(bw.MaxInt); ok {
			result = int(t)
		}
	case float32:
		result = int(t)
		if ok = t == float32(result); !ok {
			result = 0
		}
	case float64:
		result = int(t)
		if ok = t == float64(result); !ok {
			result = 0
		}
	case Number:
		val = t.val
		needRecall = true
	// case Number:
	//  val = t.val
	//  needRecall = true
	case RangeLimit:
		val = t.val
		needRecall = true
	// case *RangeLimit:
	//  val = t.val
	//  needRecall = true
	default:
		result, ok = platformSpecificInt(val)
	}
	if needRecall {
		result, ok = Int(val)
	}
	return
}

func reflectInt(val interface{}) (int, bool) {
	switch value := reflect.ValueOf(val); value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = value.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = value.Uint()
	case reflect.Float32, reflect.Float64:
		val = value.Float()
	default:
		return 0, false
	}
	return Int(val)
}

func Uint(val interface{}) (result uint, ok bool) {
	var needRecall bool
	ok = true
	switch t := val.(type) {
	case int8:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int16:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int32:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case int:
		if ok = t >= 0; ok {
			result = uint(t)
		}
	case uint8:
		result = uint(t)
	case uint16:
		result = uint(t)
	case uint32:
		result = uint(t)
	case uint:
		result = t
	case float32:
		result = uint(t)
		if ok = t == float32(result); !ok {
			result = 0
		}
	case float64:
		result = uint(t)
		if ok = t == float64(result); !ok {
			result = 0
		}
	case Number:
		val = t.val
		needRecall = true
	// case Number:
	//  val = t.val
	//  needRecall = true
	case RangeLimit:
		val = t.val
		needRecall = true
	// case *RangeLimit:
	//  val = t.val
	//  needRecall = true
	default:
		result, ok = platformSpecificUint(val)
	}
	if needRecall {
		result, ok = Uint(val)
	}
	return
}

func reflectUint(val interface{}) (result uint, ok bool) {
	switch value := reflect.ValueOf(val); value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val = value.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val = value.Uint()
	case reflect.Float32, reflect.Float64:
		val = value.Float()
	default:
		return 0, false
	}
	return Uint(val)
}

func Float64(val interface{}) (result float64, ok bool) {
	var needRecall bool
	ok = true
	switch t := val.(type) {
	case int8:
		result = float64(t)
	case int16:
		result = float64(t)
	case int32:
		result = float64(t)
	case int64:
		result = float64(t)
	case int:
		result = float64(t)
	case uint8:
		result = float64(t)
	case uint16:
		result = float64(t)
	case uint32:
		result = float64(t)
	case uint64:
		result = float64(t)
	case uint:
		result = float64(t)
	case float32:
		result = float64(t)
	case float64:
		result = t
	case Number:
		val = t.val
		needRecall = true
	// case Number:
	//  val = t.val
	//  needRecall = true
	case RangeLimit:
		val = t.val
		needRecall = true
	// case *RangeLimit:
	//  val = t.val
	//  needRecall = true
	default:
		ok = false
		needRecall = true
		switch value := reflect.ValueOf(val); value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val = value.Int()
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val = value.Uint()
		case reflect.Float32, reflect.Float64:
			val = value.Float()
		default:
			needRecall = false
		}
	}
	if needRecall {
		result, ok = Float64(val)
	}
	return
}

// ============================================================================

func MustInt(val interface{}) (result int) {
	var ok bool
	if result, ok = Int(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Int")
	}
	return
}

func MustUint(val interface{}) (result uint) {
	var ok bool
	if result, ok = Uint(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Uint")
	}
	return
}

func MustFloat64(val interface{}) (result float64) {
	var ok bool
	if result, ok = Float64(val); !ok {
		bwerr.Panic(ansiIsNotOfType, val, "Float64")
	}
	return
}

// ============================================================================

type Number struct {
	val interface{}
}

func NumberFrom(val interface{}) (result Number, ok bool) {
	var (
		i int
		u uint
		f float64
	)
	if i, ok = Int(val); ok {
		result = Number{i}
	} else if u, ok = Uint(val); ok {
		result = Number{u}
	} else if f, ok = Float64(val); ok {
		result = Number{f}
	}
	return
}

func MustNumberFrom(val interface{}) (result Number) {
	var ok bool
	if result, ok = NumberFrom(val); !ok {
		bwerr.Panic(ansi.String("<ansiVal>%#v<ansi> can not be a <ansiType>Number"), val)
	}
	return
}

func (n Number) Val() interface{} {
	return n.val
}

func (n Number) IsEqualTo(a Number) (result bool) {
	return n.compareTo(a, func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool) {
		switch kind {
		case compareUintUint:
			result = u == v
		case compareIntInt:
			result = i == j
		case compareFloat64Float64:
			result = f == g
		}
		return
	})
}

func (n Number) IsLessThan(a Number) (result bool) {
	result = n.compareTo(a, func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool) {
		switch kind {
		case compareUintUint:
			result = u < v
		case compareIntUint:
			result = true
		case compareIntInt:
			result = i < j
		case compareFloat64Float64:
			result = f < g
		}
		return
	})
	return
}

type compareKind uint8

const (
	compareUintUint compareKind = iota
	compareUintInt
	compareIntUint
	compareIntInt
	compareFloat64Float64
)

type compareFunc func(kind compareKind, u, v uint, i, j int, f, g float64) (result bool)

func (n Number) compareTo(a Number, fn compareFunc) (result bool) {
	if u, ok := Uint(n.val); ok {
		if v, ok := Uint(a.val); ok {
			result = fn(compareUintUint, u, v, 0, 0, 0, 0)
		} else if j, ok := Int(a.val); ok {
			result = fn(compareUintInt, u, 0, 0, j, 0, 0)
		} else if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, float64(u), g)
		}
	} else if i, ok := Int(n.val); ok {
		if v, ok := Uint(a.val); ok {
			result = fn(compareIntUint, 0, v, i, 0, 0, 0)
		} else if j, ok := Int(a.val); ok {
			result = fn(compareIntInt, 0, 0, i, j, 0, 0)
		} else if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, float64(i), g)
		}
	} else if f, ok := Float64(n.val); ok {
		if g, ok := Float64(a.val); ok {
			result = fn(compareFloat64Float64, 0, 0, 0, 0, f, g)
		}
	}
	return
}

func (n Number) String() string {
	bytes, _ := json.Marshal(n.val)
	return string(bytes)
}

func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.val)
}

// ============================================================================

type RangeLimit struct {
	val interface{}
}

func RangeLimitFrom(val interface{}) (result RangeLimit, ok bool) {
	var (
		n Number
	)
	if val == nil {
		result = RangeLimit{}
		ok = true
	} else if n, ok = NumberFrom(val); ok {
		result = RangeLimit{val: n.val}
	} else {
		result, ok = val.(RangeLimit)
	}
	return
}

func MustRangeLimitFrom(val interface{}) (result RangeLimit) {
	var ok bool
	if result, ok = RangeLimitFrom(val); !ok {
		bwerr.Panic(ansiValCanNotBeRangeLimit, val)
	}
	return
}

func (n RangeLimit) Val() interface{} {
	return n.val
}

func (rl RangeLimit) String() (result string) {
	if n, ok := NumberFrom(rl.val); ok {
		result = n.String()
	}
	return
}

// ============================================================================

type RangeTemplateLimit struct {
	val interface{}
}

func RangeTemplateLimitFrom(val interface{}) (result RangeTemplateLimit, ok bool) {
	var (
		path bw.ValPath
		n    Number
		rl   RangeLimit
	)
	if val == nil {
		result = RangeTemplateLimit{}
		ok = true
	} else if n, ok = NumberFrom(val); ok {
		result = RangeTemplateLimit{val: n.val}
	} else if path, ok = val.(bw.ValPath); ok {
		result = RangeTemplateLimit{val: path}
	} else if rl, ok = val.(RangeLimit); ok {
		result = RangeTemplateLimit{val: rl.val}
	} else {
		result, ok = val.(RangeTemplateLimit)
	}
	return
}

func MustRangeTemplateLimitFrom(val interface{}) (result RangeTemplateLimit) {
	var ok bool
	if result, ok = RangeTemplateLimitFrom(val); !ok {
		bwerr.Panic(ansiValCanNotBeRangeLimit, val)
	}
	return
}

func (rl RangeTemplateLimit) String() (result string) {
	if n, ok := NumberFrom(rl.val); ok {
		result = n.String()
	} else if path, ok := n.val.(bw.ValPath); !ok {
		result = path.String()
	}
	return
}

// ============================================================================

type RangeA struct {
	Min, Max interface{}
}

// ============================================================================

type Range struct {
	min, max RangeLimit
}

func (r Range) Min() RangeLimit {
	return r.min
}

func (r Range) Max() RangeLimit {
	return r.max
}

func RangeFrom(a RangeA) (result *Range, err error) {
	var min, max RangeLimit
	var ok bool
	if min, ok = RangeLimitFrom(a.Min); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Min", a.Min)
		return
	}
	if max, ok = RangeLimitFrom(a.Max); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Max", a.Max)
		return
	}
	result = &Range{min: min, max: max}
	if min, ok := NumberFrom(min.val); ok {
		if max, ok := NumberFrom(max.val); ok {
			if max.IsLessThan(min) {
				err = bwerr.From(ansiMaxMustNotBeLessThanMin, max, min)
			}
		}
	}
	return
}

func MustRangeFrom(a RangeA) (result *Range) {
	var err error
	if result, err = RangeFrom(a); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Range) String() (result string) {
	result = fmt.Sprintf("%s..%s", v.min, v.max)
	return
}

func (r Range) Contains(val interface{}) (result bool) {
	var n Number
	var ok bool
	if n, ok = NumberFrom(val); !ok {
		return false
	}
	if min, ok := NumberFrom(r.min.val); !ok {
		if n.IsEqualTo(min) {
			return true
		} else if n.IsLessThan(min) {
			return false
		}
	}
	if max, ok := NumberFrom(r.max.val); ok {
		if n.IsEqualTo(max) {
			return true
		} else if max.IsLessThan(n) {
			return false
		}
	}
	return true
}

func (r Range) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// ============================================================================

type RangeTemplate struct {
	min, max RangeTemplateLimit
}

func (r RangeTemplate) Min() RangeTemplateLimit {
	return r.min
}

func (r RangeTemplate) Max() RangeTemplateLimit {
	return r.max
}

func RangeTemplateFrom(a RangeA) (result *RangeTemplate, err error) {
	var min, max RangeTemplateLimit
	var ok bool
	if min, ok = RangeTemplateLimitFrom(a.Min); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Min", a.Min)
		return
	}
	if max, ok = RangeTemplateLimitFrom(a.Max); !ok {
		err = bwerr.From(ansiVarValCanNotBeRangeLimit, "a.Max", a.Max)
		return
	}
	result = &RangeTemplate{min: min, max: max}
	if min, ok := NumberFrom(min.val); ok {
		if max, ok := NumberFrom(max.val); ok {
			if max.IsLessThan(min) {
				err = bwerr.From(ansiMaxMustNotBeLessThanMin, max, min)
			}
		}
	}
	return
}

func MustRangeTemplateFrom(a RangeA) (result *RangeTemplate) {
	var err error
	if result, err = RangeTemplateFrom(a); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v RangeTemplate) String() (result string) {
	result = fmt.Sprintf("%s..%s", v.min, v.max)
	return
}

func (r RangeTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// ============================================================================

type EnumTemplateItem struct {
	IsPath bool
	S      string
	Path   bw.ValPathTemplate
}

// type DefTemplateEnum struct {
// 	IsPath bool
// 	Path   bw.ValPathTemplate
// 	Items  []EnumTemplateItem
// }

type DefTemplateRange struct {
	IsPath        bool
	RangeTemplate *RangeTemplate
	Path          bw.ValPathTemplate
}

type DefTemplate struct {
	Types        ValKindSet
	IsOptional   bool
	Enum         []EnumTemplateItem
	Range        *DefTemplateRange
	KeysDef      map[string]DefTemplate
	ElemDef      *DefTemplate
	ArrayElemDef *DefTemplate
	Default      *Template
	IsArrayOf    bool
}

func (v *DefTemplate) MarshalJSON() ([]byte, error) {
	return bwjson.MarshalJSON(v)
	// result := map[string]interface{}{}
	// result["Types"] = v.Types
	// result["IsOptional"] = v.IsOptional
	// if v.IsArrayOf {
	// 	result["IsArrayOf"] = v.IsArrayOf
	// }
	// if v.Enum != nil {
	// 	result["Enum"] = v.Enum
	// }
	// if v.Range != nil {
	// 	result["Range"] = v.Range
	// }
	// if v.KeysDef != nil {
	// 	result["KeysDef"] = v.KeysDef
	// }
	// if v.ElemDef != nil {
	// 	result["ElemDef"] = v.ElemDef
	// }
	// if v.ArrayElemDef != nil {
	// 	result["ArrayElemDef"] = v.ArrayElemDef
	// }
	// if v.Default != nil {
	// 	result["Default"] = v.Default
	// }
	// return json.Marshal(result)
}

// ============================================================================

type Def struct {
	Types        ValKindSet
	IsOptional   bool
	Enum         bwset.String
	Range        *Range
	KeysDef      map[string]Def
	ElemDef      *Def
	ArrayElemDef *Def
	Default      interface{}
	IsArrayOf    bool
}

func (v Def) MarshalJSON() ([]byte, error) {
	return bwjson.MarshalJSON(v)
	// result := map[string]interface{}{}
	// result["Types"] = v.Types
	// result["IsOptional"] = v.IsOptional
	// if v.IsArrayOf {
	// 	result["IsArrayOf"] = v.IsArrayOf
	// }
	// if len(v.Enum) > 0 {
	// 	result["Enum"] = v.Enum
	// }
	// if v.Range != nil {
	// 	result["Range"] = v.Range
	// }
	// if v.KeysDef != nil {
	// 	result["KeysDef"] = v.KeysDef
	// }
	// if v.ElemDef != nil {
	// 	result["ElemDef"] = v.ElemDef
	// }
	// if v.ArrayElemDef != nil {
	// 	result["ArrayElemDef"] = v.ArrayElemDef
	// }
	// if v.Default != nil {
	// 	result["Default"] = v.Default
	// }
	// return json.Marshal(result)
}

// ============================================================================

// ValKind - разновидность interface{}-значения
type ValKind uint8

// разновидности interface{}-значения
const (
	ValUnknown ValKind = iota
	ValString
	ValId
	ValBool
	ValInt
	ValUint
	ValFloat64
	ValNumber
	ValPath
	ValRange
	ValMap
	ValOrderedMap
	ValMapIntf
	ValArray
	ValDef
	ValNil
)

// MarshalJSON encoding/json support
func (v ValKind) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v ValKind) AnsiString() string {
	return ansi.String(fmt.Sprintf("<ansiType>%s<ansi>", v.String()))
}

// ============================================================================

type ValKinds []ValKind

func (vks ValKinds) Strings() (result []string) {
	for _, vk := range vks {
		result = append(result, vk.String())
	}
	return
}

// ============================================================================

//go:generate bwsetter -type=ValKind -test

//go:generate stringer -type ValKind -trimprefix Val

const (
	_ValKindSetTestItemA ValKind = ValString
	_ValKindSetTestItemB ValKind = ValBool
)

// ============================================================================

func ValKindFromString(s string) (result ValKind, err error) {
	var ok bool
	if result, ok = mapValKindFromString[s]; !ok {
		err = bwerr.From(ansiUknownValKind, result)
	}
	return
}

// ============================================================================

func Kind(val interface{}, optExpects ...ValKindSet) (result interface{}, kind ValKind) {
	expects := ValKindSet{}
	expectsAny := true
	if len(optExpects) > 0 && len(optExpects[0]) > 0 {
		expects = optExpects[0]
		expectsAny = false
	}
	if val == nil {
		if expectsAny || expects.Has(ValNil) {
			kind = ValNil
		}
	} else {
		kindOfInt := func(i int64) (result interface{}, kind ValKind) {
			if (expectsAny || expects.Has(ValInt)) && (i <= int64(bw.MaxInt)) {
				result = int(i)
				kind = ValInt
			} else if (expectsAny || expects.Has(ValUint)) && i >= 0 && (uint64(i) <= uint64(bw.MaxUint)) {
				result = uint(i)
				kind = ValUint
			} else if expectsAny || expects.Has(ValFloat64) {
				result = float64(i)
				kind = ValFloat64
			} else if expects.Has(ValNumber) {
				result = MustNumberFrom(i)
				kind = ValNumber
			}
			return
		}
		kindOfUint := func(u uint64) (result interface{}, kind ValKind) {
			if (len(expects) == 0 || expects.Has(ValInt)) && (u <= uint64(bw.MaxInt)) {
				result = int(u)
				kind = ValInt
			} else if (expectsAny || expects.Has(ValUint)) && (uint64(u) <= uint64(bw.MaxUint)) {
				result = uint(u)
				kind = ValUint
			} else if expectsAny || expects.Has(ValFloat64) {
				result = float64(u)
				kind = ValFloat64
			} else if expects.Has(ValNumber) {
				result = MustNumberFrom(u)
				kind = ValNumber
			}
			return
		}
		var needRecall bool
		switch t := val.(type) {
		case bool:
			if expectsAny || expects.Has(ValBool) {
				result = t
				kind = ValBool
			}
		case string:
			if expectsAny || expects.Has(ValString) {
				result = t
				kind = ValString
			}
		case []interface{}:
			if expectsAny || expects.Has(ValArray) {
				if !reflect.ValueOf(t).IsNil() {
					result = t
					kind = ValArray
				} else if expectsAny || expects.Has(ValNil) {
					result = nil
					kind = ValNil
				}
			}
		case map[string]interface{}:
			if expectsAny || expects.Has(ValMap) || expects.Has(ValMapIntf) {
				if !reflect.ValueOf(t).IsNil() {
					if expectsAny || expects.Has(ValMap) {
						result = t
						kind = ValMap
					} else {
						result = bwmap.M(t)
						kind = ValMapIntf
					}
				} else if expectsAny || expects.Has(ValNil) {
					result = nil
					kind = ValNil
				}
			}
		case *bwmap.Ordered:
			if expectsAny || expects.Has(ValOrderedMap) {
				result = t
				kind = ValOrderedMap
			} else if expects.Has(ValMap) {
				result = t.Map()
				kind = ValMap
			} else if expects.Has(ValMapIntf) {
				result = t
				kind = ValMapIntf
			}
		case bwmap.Ordered:
			if expectsAny || expects.Has(ValOrderedMap) {
				result = &t
				kind = ValOrderedMap
			} else if expects.Has(ValMap) {
				result = t.Map()
				kind = ValMap
			} else if expects.Has(ValMapIntf) {
				result = &t
				kind = ValMapIntf
			}
		case bwmap.I:
			if expectsAny || expects.Has(ValMapIntf) {
				result = t
				kind = ValMapIntf
			}
		case bw.ValPath:
			if expectsAny || expects.Has(ValPath) {
				result = t
				kind = ValPath
			}
		case *Range:
			if expectsAny || expects.Has(ValRange) {
				result = t
				kind = ValRange
			}
		case Range:
			if expectsAny || expects.Has(ValRange) {
				result = &t
				kind = ValRange
			}
		case *Def:
			if expectsAny || expects.Has(ValDef) {
				result = t
				kind = ValDef
			}
		case Def:
			if expectsAny || expects.Has(ValDef) {
				result = &t
				kind = ValDef
			}
		case int:
			result, kind = kindOfInt(int64(t))
		case int8:
			result, kind = kindOfInt(int64(t))
		case int16:
			result, kind = kindOfInt(int64(t))
		case int32:
			result, kind = kindOfInt(int64(t))
		case int64:
			result, kind = kindOfInt(t)
		case uint:
			result, kind = kindOfUint(uint64(t))
		case uint8:
			result, kind = kindOfUint(uint64(t))
		case uint16:
			result, kind = kindOfUint(uint64(t))
		case uint64:
			result, kind = kindOfUint(t)
		case float32:
			if (expectsAny || expects.Has(ValInt)) && float32(int(t)) == t {
				result = int(t)
				kind = ValInt
			} else if (expectsAny || expects.Has(ValUint)) && float32(uint(t)) == t {
				result = uint(t)
				kind = ValUint
			} else if expectsAny || expects.Has(ValFloat64) {
				result = float64(t)
				kind = ValFloat64
			} else if expects.Has(ValNumber) {
				result = MustNumberFrom(t)
				kind = ValNumber
			}
		case float64:
			if (expectsAny || expects.Has(ValInt)) && (float64(int(t)) == t) {
				result = int(t)
				kind = ValInt
			} else if (expectsAny || expects.Has(ValUint)) && (float64(uint(t)) == t) {
				result = uint(t)
				kind = ValUint
			} else if expectsAny || expects.Has(ValFloat64) {
				result = t
				kind = ValFloat64
			} else if expects.Has(ValNumber) {
				result = MustNumberFrom(t)
				kind = ValNumber
			}
		case Number:
			if expects.Has(ValNumber) {
				result = t
				kind = ValNumber
			} else {
				val = t.Val()
				needRecall = true
			}
		case RangeLimit:
			val = t.Val()
			needRecall = true
		default:
			switch value := reflect.ValueOf(val); value.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if expectsAny || expects.HasAny(ValInt, ValUint, ValFloat64) {
					val = value.Int()
					needRecall = true
				}
			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				if expectsAny || expects.HasAny(ValInt, ValUint, ValFloat64) {
					val = value.Uint()
					needRecall = true
				}
			case reflect.Float32, reflect.Float64:
				if expectsAny || expects.HasAny(ValInt, ValUint, ValFloat64) {
					val = value.Float()
					needRecall = true
				}
			case reflect.Bool:
				if expectsAny || expects.Has(ValBool) {
					result = value.Bool()
					kind = ValBool
				}
			case reflect.String:
				if expectsAny || expects.Has(ValString) {
					result = value.String()
					kind = ValString
				}
			case reflect.Slice:
				if expectsAny || expects.Has(ValArray) {
					if !reflect.ValueOf(t).IsNil() {
						len := value.Len()
						vals := []interface{}{}
						for i := 0; i < len; i++ {
							vals = append(vals, value.Index(i).Interface())
						}
						result = vals
						kind = ValArray
					} else if expectsAny || expects.Has(ValNil) {
						result = nil
						kind = ValNil
					}
				}
			case reflect.Map:
				if expectsAny || expects.Has(ValMap) || expects.Has(ValMapIntf) {
					if reflect.TypeOf(val).Key().Kind() == reflect.String {
						if !reflect.ValueOf(t).IsNil() {
							m := map[string]interface{}{}
							keyValues := value.MapKeys()
							len := len(keyValues)
							for i := 0; i < len; i++ {
								keyValue := keyValues[i]
								m[keyValue.String()] = value.MapIndex(keyValue).Interface()
							}
							if expectsAny || expects.Has(ValMap) {
								result = m
								kind = ValMap
							} else {
								result = bwmap.M(m)
								kind = ValMapIntf
							}
						} else if expectsAny || expects.Has(ValNil) {
							result = nil
							kind = ValNil
						}
					}
				}
			}
		}
		if needRecall {
			result, kind = Kind(val, expects)
		}
	}
	return
}

// ============================================================================

var (
	mapValKindFromString = map[string]ValKind{}

	ansiUknownValKind            string
	ansiVarValCanNotBeRangeLimit string
	ansiValCanNotBeRangeLimit    string
	// ansiIsNotOfType              string
	ansiMaxMustNotBeLessThanMin string
)

func init() {
	for i := ValUnknown; i <= ValNil; i++ {
		mapValKindFromString[i.String()] = i
	}
	ansiUknownValKind = ansi.String("<ansiPath>ValKindFromString<ansi>: uknown <ansiVal>%s")
	// ansiVarValCanNotBeRangeLimit = ansi.String("<ansiVar>%s<ansi> (<ansiVal>%#v<ansi>) can not be a <ansiType>RangeLimit")
	// ansiValCanNotBeRangeLimit = ansi.String("<ansiVal>%#v<ansi> can not be a <ansiType>RangeLimit")
	// ansiIsNotOfType = ansi.String("<ansiVal>%#v<ansi> is not <ansiType>%s")
	ansiMaxMustNotBeLessThanMin = ansi.String("<ansiVar>a.Max<ansi> (<ansiVal>%s<ansi>) must not be <ansiErr>less<ansi> then <ansiVar>a.Min<ansi> (<ansiVal>%s<ansi>)")
}

// ============================================================================
