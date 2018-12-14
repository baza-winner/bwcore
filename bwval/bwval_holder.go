package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

type Holder struct {
	Val interface{}
	Pth bw.ValPath
}

// ============================================================================

// PathVal - реализация интерфейса bw.Val
func (v Holder) PathVal(path bw.ValPath, optVars ...map[string]interface{}) (result interface{}, err error) {
	if len(path) == 0 {
		result = v.Val
		return
	}
	defer func() {
		if err != nil {
			result = nil
		}
	}()

	var simplePath bw.ValPath
	simplePath, err = v.simplifyPath(path, optVars)
	if err != nil {
		return
	}

	if path[0].Type == bw.ValPathItemVar {
		varName := path[0].Key
		var target interface{}
		var ok bool
		if ok = len(optVars) > 0; ok {
			target, ok = optVars[0][varName]
		}
		if !ok && !hasOptional(path) {
			err = bwerr.From(ansi.String("var <ansiVar>%s<ansi> is not defined"), varName)
			return
		}
		h := Holder{Val: target}
		return h.PathVal(simplePath[1:])
	}

	result = v.Val
	for i, vpi := range simplePath {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = Holder{result, path[:i]}.KeyVal(vpi.Key,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
		case bw.ValPathItemIdx:
			result, err = Holder{result, path[:i]}.IdxVal(vpi.Idx,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
		case bw.ValPathItemHash:
			if result == nil {
				result = 0
			} else {
				switch t := result.(type) {
				case map[string]interface{}:
					result = len(t)
				case []interface{}:
					result = len(t)
				default:
					// err = Holder{result, path[:i]}.notOfValKindError("Map", "Array")
					err = Holder{result, path[:i]}.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValMap, bwtype.ValArray))
				}
			}
		}
		if err != nil {
			return
		}
	}
	return
}

func (v Holder) Path(pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result Holder, err error) {
	var val interface{}
	var path bw.ValPath
	if path, err = pathProvider.Path(); err != nil {
		err = bwerr.Refine(err, "invalid path: {Error}")
		return
	}
	if val, err = (&v).PathVal(path, optVars...); err != nil {
		return
	}
	result = Holder{val, path}
	return
}

func (v Holder) MustPath(pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result Holder) {
	var err error
	if result, err = v.Path(pathProvider, optVars...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) MustPathStr(pathStr string, optVars ...map[string]interface{}) (result Holder) {
	return v.MustPath(PathS{S: pathStr})
}

// MarshalJSON - реализация интерфейса bw.Val
func (v Holder) MarshalJSON() ([]byte, error) {
	if len(v.Pth) == 0 {
		return json.Marshal(v.Val)
	} else {
		result := map[string]interface{}{}
		result["val"] = v.Val
		result["path"] = v.Pth.String()
		return json.Marshal(result)
	}
}

// SetPathVal - реализация интерфейса bw.Val
func (v *Holder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	if len(path) == 0 {
		v.Val = val
		return
	}
	if path[len(path)-1].Type == bw.ValPathItemHash {
		return readonlyPathError(path)
	}

	var simplePath bw.ValPath
	simplePath, err = v.simplifyPath(path, optVars)
	if err != nil {
		return
	}

	result := v.Val
	if result == nil {
		return v.wrongValError()
	}

	if path[0].Type == bw.ValPathItemVar {
		var vars map[string]interface{}
		if len(optVars) > 0 {
			vars = optVars[0]
		}
		if vars == nil {
			return bwerr.From(ansiVarsIsNil)
		}
		simplePath[0].Type = bw.ValPathItemKey
		h := Holder{Val: vars}
		return h.SetPathVal(val, simplePath)
	}

	if len(simplePath) > 1 {
		for i, vpi := range simplePath[:len(simplePath)-1] {
			switch vpi.Type {
			case bw.ValPathItemKey:
				hVal := Holder{result, path[:i+1]}
				result, err = hVal.KeyVal(
					vpi.Key,
					func() (result interface{}, ok bool) {
						if ok = simplePath[i+1].Type == bw.ValPathItemKey; ok {
							_, _ = hVal.KindSwitch(map[bwtype.ValKind]KindCase{
								bwtype.ValMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
									m, _ := val.(map[string]interface{})
									result = map[string]interface{}{}
									m[vpi.Key] = result
									return nil, nil
								},
								bwtype.ValOrderedMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
									o, _ := val.(*bwmap.Ordered)
									result = bwmap.OrderedNew()
									o.Set(vpi.Key, result)
									return nil, nil
								},
							})
						}
						return
					},
				)
			case bw.ValPathItemIdx:
				result, err = Holder{result, path[:i+1]}.IdxVal(vpi.Idx)
			}
			if err != nil {
				return
			} else if result == nil {
				return Holder{nil, path[:i+1]}.wrongValError()
			}
		}
	}
	rh := Holder{result, path[:len(path)-1]}
	vpi := simplePath[len(simplePath)-1]
	switch vpi.Type {
	case bw.ValPathItemKey:
		// bwdebug.Print("!LOOK AT ME", "vpi.Key", vpi.Key, "val", val)
		err = rh.SetKeyVal(vpi.Key, val)
	case bw.ValPathItemIdx:
		err = rh.SetIdxVal(vpi.Idx, val)
	}
	return
}

// ============================================================================

func (v Holder) Bool() (result bool, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValBool)); kind == bwtype.ValBool {
		result, _ = t.(bool)
		ok = true
	}
	return
}

func (v Holder) MustBool(optDefault ...bool) (result bool) {
	var ok bool
	if result, ok = v.Bool(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValBool)))
		}
	}
	return
}

func (v Holder) String() (result string, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValString)); kind == bwtype.ValString {
		result, _ = t.(string)
		ok = true
	}
	return
}

func (v Holder) MustString(optDefault ...string) (result string) {
	var ok bool
	if result, ok = v.String(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValString)))
		}
	}
	return
}

func (v Holder) Int() (result int, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValInt)); kind == bwtype.ValInt {
		result, _ = t.(int)
		ok = true
	}
	return
}

func (v Holder) MustInt(optDefault ...int) (result int) {
	var ok bool
	if result, ok = v.Int(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValInt)))
		}
	}
	return
}

func (v Holder) Uint() (result uint, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValUint)); kind == bwtype.ValUint {
		result, _ = t.(uint)
		ok = true
	}
	return
}

func (v Holder) MustUint(optDefault ...uint) (result uint) {
	var ok bool
	if result, ok = v.Uint(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValUint)))
		}
	}
	return
}

func (v Holder) Float64() (result float64, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValFloat64)); kind == bwtype.ValFloat64 {
		result, _ = t.(float64)
		ok = true
	}
	return
}

func (v Holder) MustFloat64(optDefault ...float64) (result float64) {
	var ok bool
	if result, ok = v.Float64(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValFloat64)))
		}
	}
	return
}

func (v Holder) Array() (result []interface{}, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValArray)); kind == bwtype.ValArray {
		v.Val = t
		result, _ = t.([]interface{})
		if len(result) == 0 {
			result = []interface{}{}
		}
		ok = true
	}
	return
}

func (v Holder) MustArray(optDefault ...[]interface{}) (result []interface{}) {
	var ok bool
	if result, ok = v.Array(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValArray)))
		}
	}
	return
}

func (v Holder) Map() (result map[string]interface{}, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValMap)); kind == bwtype.ValMap {
		v.Val = t
		result, _ = t.(map[string]interface{})
		ok = true
	}
	return
}

func (v Holder) MustMap(optDefault ...map[string]interface{}) (result map[string]interface{}) {
	var ok bool
	if result, ok = v.Map(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValMap)))
		}
	}
	return
}

func (v Holder) OrderedMap() (result *bwmap.Ordered, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValOrderedMap)); kind == bwtype.ValOrderedMap {
		result, _ = t.(*bwmap.Ordered)
		ok = true
	}
	return
}

func (v Holder) MustOrderedMap(optDefault ...*bwmap.Ordered) (result *bwmap.Ordered) {
	var ok bool
	if result, ok = v.OrderedMap(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValOrderedMap)))
		}
	}
	return
}

func (v Holder) MapIntf() (result bwmap.I, ok bool) {
	if t, kind := bwtype.Kind(v.Val, bwtype.ValKindSetFrom(bwtype.ValMapIntf)); kind == bwtype.ValMapIntf {
		result, _ = t.(bwmap.I)
		ok = true
	}
	return
}

func (v Holder) MustMapIntf(optDefault ...bwmap.I) (result bwmap.I) {
	var ok bool
	if result, ok = v.MapIntf(); !ok {
		if len(optDefault) > 0 {
			result = optDefault[0]
		} else {
			bwerr.PanicErr(v.notOfValKindError(bwtype.ValKindSetFrom(bwtype.ValMapIntf)))
		}
	}
	return
}

// ============================================================================

type ForEachBody func(idx int, key string, hVal Holder) (needBreak bool, err error)

func (v Holder) ForEach(body ForEachBody) (err error) {
	var needBreak bool
	_, err = v.KindSwitch(map[bwtype.ValKind]KindCase{
		bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
			o, _ := val.(bwmap.I)
			for idx, key := range o.Keys() {
				hVal := v.MustKey(key)
				if needBreak, err = body(idx, key, hVal); needBreak || err != nil {
					break
				}
			}
			return
		},
		// bwtype.ValOrderedMap: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
		// 	o, _ := val.(*bwmap.Ordered)
		// 	for idx, key := range o.Keys() {
		// 		hVal := v.MustKey(key)
		// 		if needBreak, err = body(idx, key, hVal); needBreak || err != nil {
		// 			break
		// 		}
		// 	}
		// 	return
		// },
		// bwtype.ValMap: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
		// 	idx := 0
		// 	m, _ := val.(map[string]interface{})
		// 	for key, _ := range m {
		// 		hVal := v.MustKey(key)
		// 		if needBreak, err = body(idx, key, hVal); needBreak || err != nil {
		// 			break
		// 		}
		// 		idx++
		// 	}
		// 	return
		// },
		bwtype.ValArray: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
			vals, _ := val.([]interface{})
			for idx, _ := range vals {
				hVal := v.MustIdx(idx)
				if needBreak, err = body(idx, "", hVal); needBreak || err != nil {
					break
				}
			}
			return
		},
	})
	return
}

// ============================================================================

func (v Holder) Keys(optFilter ...bwmap.KeysFilter) (result []string, err error) {
	_, _ = v.KindSwitch(map[bwtype.ValKind]KindCase{
		bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			m, _ := val.(bwmap.I)
			result = m.Keys(optFilter...)
			return nil, nil
		},
		// bwtype.ValOrderedMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	o, _ := val.(*bwmap.Ordered)
		// 	result = o.Keys(optFilter...)
		// 	return nil, nil
		// },
	})
	return
}

func (v Holder) MustKeys(optFilter ...bwmap.KeysFilter) (result []string) {
	var err error
	if result, err = v.Keys(optFilter...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) HasKey(key string) (result bool) {
	_, _ = v.KindSwitch(map[bwtype.ValKind]KindCase{
		// bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	m, _ := val.(bwmap.I)
		// 	_, result = m[key]
		// 	return nil, nil
		// },
		bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			o, _ := val.(bwmap.I)
			result = o.HasKey(key)
			return nil, nil
		},
	})
	return
}

func (v Holder) DelKey(key string) (result bool) {
	_, _ = v.KindSwitch(map[bwtype.ValKind]KindCase{
		// bwtype.ValMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	m, _ := val.(map[string]interface{})
		// 	_, result = m[key]
		// 	return nil, nil
		// },
		bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			o, _ := val.(bwmap.I)
			o.DelKey(key)
			return nil, nil
		},
	})
	return
}

func (v Holder) KeyVal(key string, optDefaultValProvider ...defaultValProvider) (result interface{}, err error) {
	// if v.Val == nil {
	// 	var ok bool
	// 	if result, ok = defaultVal(optDefaultValProvider); !ok {
	// 		err = v.wrongValError()
	// 	}
	// 	return
	// }
	_, err = v.KindSwitch(map[bwtype.ValKind]KindCase{
		bwtype.ValNil: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			var (
				ok  bool
				err error
			)
			if result, ok = defaultVal(optDefaultValProvider); !ok {
				err = v.wrongValError()
			}
			return val, err
		},
		// bwtype.ValMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	m, _ := val.(map[string]interface{})
		// 	var (
		// 		ok  bool
		// 		err error
		// 	)
		// 	if result, ok = m[key]; !ok {
		// 		if result, ok = defaultVal(optDefaultValProvider); !ok {
		// 			err = v.hasNoKeyError(key)
		// 		}
		// 	}
		// 	return val, err
		// },
		// bwtype.ValOrderedMap: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	o, _ := val.(*bwmap.Ordered)
		// 	var (
		// 		ok  bool
		// 		err error
		// 	)
		// 	if result, ok = o.Get(key); !ok {
		// 		if result, ok = defaultVal(optDefaultValProvider); !ok {
		// 			err = v.hasNoKeyError(key)
		// 		}
		// 	}
		// 	return val, err
		// },
		bwtype.ValMapIntf: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			o, _ := val.(bwmap.I)
			var (
				ok  bool
				err error
			)
			if result, ok = o.Get(key); !ok {
				if result, ok = defaultVal(optDefaultValProvider); !ok {
					err = v.hasNoKeyError(key)
				}
			}
			return val, err
		},
	})
	return
}

func (v Holder) MustKeyVal(key string, optDefaultValProvider ...defaultValProvider) (result interface{}) {
	var err error
	if result, err = v.KeyVal(key, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Key(key string, optDefaultValProvider ...defaultValProvider) (result Holder, err error) {
	var val interface{}
	if val, err = v.KeyVal(key, optDefaultValProvider...); err == nil {
		result = Holder{val, v.Pth.AppendKey(key)}
	}
	return
}

func (v Holder) MustKey(key string, optDefaultValProvider ...defaultValProvider) (result Holder) {
	var err error
	if result, err = v.Key(key, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v *Holder) SetKeyVal(key string, val interface{}) (err error) {
	_, err = v.KindSwitch(map[bwtype.ValKind]KindCase{
		// bwtype.ValMap: func(valMap interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	m, _ := valMap.(map[string]interface{})
		// 	m[key] = val
		// 	return nil, nil
		// },
		// bwtype.ValOrderedMap: func(valOrdered interface{}, kind bwtype.ValKind) (interface{}, error) {
		// 	o, _ := valOrdered.(*bwmap.Ordered)
		// 	o.Set(key, val)
		// 	return nil, nil
		// },
		bwtype.ValMapIntf: func(valOrdered interface{}, kind bwtype.ValKind) (interface{}, error) {
			o, _ := valOrdered.(bwmap.I)
			o.Set(key, val)
			return nil, nil
		},
	})
	return
}

func (v Holder) MustSetKeyVal(key string, val interface{}) {
	var err error
	if err = v.SetKeyVal(key, val); err != nil {
		bwerr.PanicErr(err)
	}
}

// ============================================================================

func (v Holder) HasIdx(idx int) (result bool) {
	if v.Val == nil {
		return
	}
	_ = v.idxHelper(idx,
		func(vals []interface{}, nidx int, ok bool) (err error) {
			result = ok
			return
		},
	)
	return
}

func (v Holder) IdxVal(idx int, optDefaultValProvider ...defaultValProvider) (result interface{}, err error) {
	if v.Val == nil {
		var ok bool
		if len(optDefaultValProvider) > 0 {
			result, ok = optDefaultValProvider[0]()
		}
		if !ok {
			err = v.wrongValError()
		}
		return
	}
	err = v.idxHelper(idx,
		func(vals []interface{}, nidx int, ok bool) (err error) {
			if ok {
				result = vals[nidx]
			} else {
				result, ok = defaultVal(optDefaultValProvider)
			}
			if !ok {
				err = v.notEnoughRangeError(len(vals), idx)
			}
			return
		},
	)
	return
}

func (v Holder) MustIdxVal(idx int, optDefaultValProvider ...defaultValProvider) (result interface{}) {
	var err error
	if result, err = v.IdxVal(idx, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Idx(idx int, optDefaultValProvider ...defaultValProvider) (result Holder, err error) {
	var val interface{}
	if val, err = v.IdxVal(idx, optDefaultValProvider...); err == nil {
		result = Holder{val, v.Pth.AppendIdx(idx)}
	}
	return
}

func (v Holder) MustIdx(idx int, optDefaultValProvider ...defaultValProvider) (result Holder) {
	var err error
	if result, err = v.Idx(idx, optDefaultValProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) SetIdxVal(idx int, val interface{}) (err error) {
	err = v.idxHelper(idx,
		func(vals []interface{}, nidx int, ok bool) (err error) {
			if !ok {
				err = v.notEnoughRangeError(len(vals), idx)
			} else {
				vals[nidx] = val
			}
			return
		},
	)
	return
}

// ============================================================================

func (v Holder) ValidVal(def Def) (result interface{}, err error) {
	result, err = v.validVal(def, false)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	}
	return
}

func (v Holder) MustValidVal(def Def) (result interface{}) {
	var err error
	if result, err = v.ValidVal(def); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v *Holder) Valid(def Def) (result Holder, err error) {
	var val interface{}
	val, err = v.validVal(def, false)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	} else {
		result = Holder{Val: val, Pth: v.Pth}
	}
	return
}

func (v Holder) MustValid(def Def) (result Holder) {
	var err error
	if result, err = v.Valid(def); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

type KindCase func(val interface{}, kind bwtype.ValKind) (interface{}, error)

func (v Holder) KindSwitch(kindCases map[bwtype.ValKind]KindCase, optDefaultCase ...KindCase) (val interface{}, err error) {
	expects := bwtype.ValKindSet{}
	for k, _ := range kindCases {
		expects.Add(k)
	}
	val, kind := bwtype.Kind(v.Val, expects)
	if KindCase, ok := kindCases[kind]; ok {
		val, err = KindCase(val, kind)
	} else if len(optDefaultCase) == 0 {
		vkSet := bwtype.ValKindSet{}
		for vk := range kindCases {
			vkSet.Add(vk)
		}
		err = v.notOfValKindError(vkSet)
	} else if optDefaultCase[0] != nil {
		val, err = optDefaultCase[0](val, kind)
	}
	return
}

func (v Holder) MustKindSwitch(kindCases map[bwtype.ValKind]KindCase, optDefaultCase ...KindCase) (val interface{}) {
	var err error
	if val, err = v.KindSwitch(kindCases, optDefaultCase...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

func (v Holder) validVal(def Def, skipArrayOf bool) (result interface{}, err error) {
	var (
		defKind    bwtype.ValKind
		val        interface{}
		needReturn bool
	)
	setKind := func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
		if def.Types.Has(kind) {
			defKind = kind
		}
		return val, nil
	}
	// // bwdebug.Print("v.Val:#v", v.Val)
	if val, err = v.KindSwitch(map[bwtype.ValKind]KindCase{
		bwtype.ValNil: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
			result = val
			if !skipArrayOf {
				if def.Default != nil {
					result = def.Default
					needReturn = true
					return
				}
				if def.IsOptional {
					needReturn = true
					return
				}
			}
			if def.Types.Has(bwtype.ValMap) {
				result = map[string]interface{}{}
				defKind = bwtype.ValMap
			} else {
				err = v.notOfValKindError(def.Types)
			}
			// // bwdebug.Print("result", result, "defKind", defKind)
			return
		},
		bwtype.ValBool: setKind,
		bwtype.ValMap:  setKind,
		bwtype.ValOrderedMap: func(val interface{}, kind bwtype.ValKind) (result interface{}, err error) {
			// // bwdebug.Print("!HERE")
			if def.Types.Has(kind) {
				defKind = kind
				result = val
			} else if def.Types.Has(bwtype.ValMap) {
				defKind = bwtype.ValMap
				o, _ := val.(*bwmap.Ordered)
				result = o.Map()
			}
			return
		},
		bwtype.ValString: setKind,
		bwtype.ValInt: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			if def.Types.Has(bwtype.ValInt) {
				defKind = bwtype.ValInt
			} else if def.Types.Has(bwtype.ValNumber) {
				defKind = bwtype.ValNumber
			}
			return val, nil
		},
		bwtype.ValFloat64: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			if def.Types.Has(bwtype.ValNumber) {
				defKind = bwtype.ValNumber
			}
			return val, nil
		},
		bwtype.ValArray: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			if def.Types.Has(bwtype.ValArray) || !skipArrayOf && def.IsArrayOf {
				defKind = bwtype.ValArray
			}
			return val, nil
		},
	}, nil); err != nil || needReturn {
		result = val
		return
	}

	if defKind == bwtype.ValUnknown {
		err = v.notOfValKindError(def.Types)
		return
	}

	switch defKind {
	case bwtype.ValBool:
	case bwtype.ValString:
		if def.Enum != nil {
			s, _ := val.(string)
			if !def.Enum.Has(s) {
				err = v.unexpectedEnumValueError(def.Enum)
				return
			}
		}
	case bwtype.ValInt, bwtype.ValNumber:
		if !def.Range.Contains(val) {
			err = v.outOfRangeError(def.Range)
			return
		}

	case bwtype.ValMap, bwtype.ValOrderedMap:
		var m map[string]interface{}
		if defKind == bwtype.ValMap {
			m, _ = val.(map[string]interface{})
		} else {
			o, _ := val.(bwmap.Ordered)
			m = o.Map()
		}
		if def.Keys != nil {
			unexpectedKeys := bwmap.MustUnexpectedKeys(val, def.Keys)
			for key, keyDef := range def.Keys {
				if err = mapHelper(v.Pth, m, key, keyDef); err != nil {
					return
				}
			}
			if unexpectedKeys != nil {
				if def.Elem == nil {
					err = v.unexpectedKeysError(unexpectedKeys)
					return
				} else {
					for _, key := range unexpectedKeys.ToSlice() {
						if err = mapHelper(v.Pth, m, key, *(def.Elem)); err != nil {
							return
						}
					}
				}
			}
		} else if def.Elem != nil {
			for k := range m {
				if err = mapHelper(v.Pth, m, k, *(def.Elem)); err != nil {
					return
				}
			}
		}
	case bwtype.ValArray:
		if !skipArrayOf && def.IsArrayOf {
			v.Val = val
			if val, err = v.arrayHelper(def, true); err != nil {
				return
			}
		} else {
			elemDef := def.ArrayElem
			if elemDef == nil {
				elemDef = def.Elem
			}
			if elemDef != nil {
				v.Val = val
				if val, err = v.arrayHelper(*elemDef, false); err != nil {
					return
				}
			}
		}
	}

	if !skipArrayOf && def.IsArrayOf && defKind != bwtype.ValArray {
		val = []interface{}{val}
	}

	result = val
	return
}

func mapHelper(path bw.ValPath, m map[string]interface{}, key string, elemDef Def) (err error) {
	vp, _ := Holder{Val: m, Pth: path}.Key(key)
	var val interface{}
	if val, err = vp.validVal(elemDef, false); err != nil {
		return
	} else if val != nil {
		m[key] = val
		// // bwdebug.Print("m", m)
	}
	return
}

func (v Holder) arrayHelper(elemDef Def, skipArrayOf bool) (result interface{}, err error) {
	arr := v.MustArray()
	newArr := make([]interface{}, 0, len(arr))
	for i := range arr {
		var vp Holder
		if vp, err = v.Idx(i); err != nil {
			break
		} else {
			var val interface{}
			if val, err = vp.validVal(elemDef, skipArrayOf); err != nil {
				return
			}
			newArr = append(newArr, val)
		}
	}
	return newArr, err
}

// ============================================================================

func (v Holder) simplifyPath(path bw.ValPath, optVars []map[string]interface{}) (result bw.ValPath, err error) {
	result = bw.ValPath{}
	for _, vpi := range path {
		if vpi.Type != bw.ValPathItemPath {
			result = append(result, vpi)
		} else {
			var val interface{}
			if val, err = v.PathVal(vpi.Path, optVars...); err != nil {
				return
			}
			h := Holder{Val: val}
			if _, err = h.KindSwitch(map[bwtype.ValKind]KindCase{
				bwtype.ValString: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
					s, _ := val.(string)
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: s})
					return val, nil
				},
				bwtype.ValInt: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
					i, _ := val.(int)
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: i})
					return val, nil
				},
			}); err != nil {
				return
			}
		}
	}
	return
}

func (v Holder) idxHelper(
	idx int,
	onArray func(vals []interface{}, nidx int, ok bool) error,
) (err error) {
	var nidx int
	var ok bool
	_, err = v.KindSwitch(map[bwtype.ValKind]KindCase{
		bwtype.ValArray: func(val interface{}, kind bwtype.ValKind) (interface{}, error) {
			vals, _ := val.([]interface{})
			nidx, ok = bw.NormalIdx(idx, len(vals))
			return vals, onArray(vals, nidx, ok)
		},
	})
	return
}

type defaultValProvider func() (interface{}, bool)

func defaultVal(optDefaultValProvider []defaultValProvider) (result interface{}, ok bool) {
	if len(optDefaultValProvider) > 0 {
		if optDefaultValProvider[0] == nil {
			result = nil
			ok = true
		} else {
			result, ok = optDefaultValProvider[0]()
		}
	}
	return
}

func hasOptional(path bw.ValPath) bool {
	for _, vpi := range path {
		if vpi.IsOptional {
			return true
		}
	}
	return false
}

// ============================================================================
