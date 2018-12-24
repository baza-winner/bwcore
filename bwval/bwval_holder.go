package bwval

import (
	"encoding/json"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	bwjson "github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwmap"
)

// ============================================================================

type Holder struct {
	Val interface{}
	Pth bw.ValPath
	Var map[string]interface{}
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
		varName := path[0].Name
		var target interface{}
		var ok bool
		if ok = len(optVars) > 0; ok {
			target, ok = optVars[0][varName]
		}
		if !ok && !hasOptional(path) {
			err = bwerr.From(ansi.String("var <ansiVar>%s<ansi> is not defined"), varName)
			return
		}
		h := Holder{Val: target, Pth: bw.ValPath{path[0]}}
		return h.PathVal(simplePath[1:])
	}

	result = v.Val
	for i, vpi := range simplePath {
		switch vpi.Type {
		case bw.ValPathItemKey:
			result, err = Holder{Val: result, Pth: v.Pth.Append(path[:i])}.KeyVal(vpi.Name,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
		case bw.ValPathItemIdx:
			result, err = Holder{Val: result, Pth: v.Pth.Append(path[:i])}.IdxVal(vpi.Idx,
				func() (result interface{}, ok bool) {
					ok = hasOptional(path[i:])
					return
				},
			)
			// case bw.ValPathItemHash:
			// 	if result == nil {
			// 		result = 0
			// 	} else {
			// 		switch t := result.(type) {
			// 		case map[string]interface{}:
			// 			result = len(t)
			// 		case []interface{}:
			// 			result = len(t)
			// 		default:
			// 			// err = Holder{result, path[:i]}.notOfValKindError("Map", "Array")
			// 			err = Holder{result, v.Pth.Append(path[:i])}.notOfValKindError(ValKindSetFrom(ValMap, ValArray))
			// 		}
			// 	}
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
	result = Holder{Val: val, Pth: path}
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
	// if len(v.Pth) == 0 {
	// 	return json.Marshal(v.Val)
	// } else {
	result := map[string]interface{}{}
	result["val"] = v.Val
	result["path"] = v.Pth.String()
	return json.Marshal(result)
	// }
}

// SetPathVal - реализация интерфейса bw.Val
func (v *Holder) SetPathVal(val interface{}, path bw.ValPath, optVars ...map[string]interface{}) (err error) {
	if len(path) == 0 {
		v.Val = val
		return
	}
	// if path[len(path)-1].Type == bw.ValPathItemHash {
	// 	return readonlyPathError(path)
	// }

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
				hVal := Holder{Val: result, Pth: path[:i+1]}
				result, err = hVal.KeyVal(
					vpi.Name,
					func() (result interface{}, ok bool) {
						if ok = simplePath[i+1].Type == bw.ValPathItemKey; ok {
							_, _ = hVal.KindSwitch(map[ValKind]KindCase{
								ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
									o, _ := val.(bwmap.I)
									result = bwmap.OrderedNew()
									o.Set(vpi.Name, result)
									return nil, nil
								},
								// ValMap: func(val interface{}, kind ValKind) (interface{}, error) {
								// 	m, _ := val.(map[string]interface{})
								// 	result = map[string]interface{}{}
								// 	m[vpi.Key] = result
								// 	return nil, nil
								// },
								// ValOrderedMap: func(val interface{}, kind ValKind) (interface{}, error) {
								// 	o, _ := val.(*bwmap.Ordered)
								// 	result = bwmap.OrderedNew()
								// 	o.Set(vpi.Key, result)
								// 	return nil, nil
								// },
							})
						}
						return
					},
				)
			case bw.ValPathItemIdx:
				result, err = Holder{Val: result, Pth: v.Pth.Append(path[:i+1])}.IdxVal(vpi.Idx)
			}
			if err != nil {
				return
			} else if result == nil {
				return Holder{Pth: v.Pth.Append(path[:i+1])}.wrongValError()
			}
		}
	}
	rh := Holder{Val: result, Pth: path[:len(path)-1]}
	vpi := simplePath[len(simplePath)-1]
	switch vpi.Type {
	case bw.ValPathItemKey:
		// bwdebug.Print("!LOOK AT ME", "vpi.Key", vpi.Key, "val", val)
		err = rh.SetKeyVal(vpi.Name, val)
	case bw.ValPathItemIdx:
		err = rh.SetIdxVal(vpi.Idx, val)
	}
	return
}

// ============================================================================

func (v Holder) Bool(optDefaultProvider ...func() bool) (result bool, err error) {
	expectedType := ValBool
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(bool)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustBool(optDefaultProvider ...func() bool) (result bool) {
	var err error
	if result, err = v.Bool(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) String(optDefaultProvider ...func() string) (result string, err error) {
	expectedType := ValString
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(string)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
		ValNumber: func(val interface{}, kind ValKind) (interface{}, error) {
			n, _ := val.(Number)
			result = n.String()
			return nil, nil
		},
	})
	return
}

func (v Holder) MustString(optDefaultProvider ...func() string) (result string) {
	var err error
	if result, err = v.String(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Int(optDefaultProvider ...func() int) (result int, err error) {
	expectedType := ValInt
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(int)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustInt(optDefaultProvider ...func() int) (result int) {
	var err error
	if result, err = v.Int(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Uint(optDefaultProvider ...func() uint) (result uint, err error) {
	expectedType := ValUint
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(uint)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustUint(optDefaultProvider ...func() uint) (result uint) {
	var err error
	if result, err = v.Uint(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Float64(optDefaultProvider ...func() float64) (result float64, err error) {
	expectedType := ValFloat64
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(float64)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustFloat64(optDefaultProvider ...func() float64) (result float64) {
	var err error
	if result, err = v.Float64(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Array(optDefaultProvider ...func() []interface{}) (result []interface{}, err error) {
	result = []interface{}{}
	expectedType := ValArray
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.([]interface{})
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustArray(optDefaultProvider ...func() []interface{}) (result []interface{}) {
	var err error
	if result, err = v.Array(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) Map(optDefaultProvider ...func() map[string]interface{}) (result map[string]interface{}, err error) {
	result = map[string]interface{}{}
	expectedType := ValMap
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(map[string]interface{})
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustMap(optDefaultProvider ...func() map[string]interface{}) (result map[string]interface{}) {
	var err error
	if result, err = v.Map(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) OrderedMap(optDefaultProvider ...func() *bwmap.Ordered) (result *bwmap.Ordered, err error) {
	result = bwmap.OrderedNew()
	expectedType := ValOrderedMap
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(*bwmap.Ordered)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustOrderedMap(optDefaultProvider ...func() *bwmap.Ordered) (result *bwmap.Ordered) {
	var err error
	if result, err = v.OrderedMap(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v Holder) MapIntf(optDefaultProvider ...func() bwmap.I) (result bwmap.I, err error) {
	result = bwmap.M(map[string]interface{}{})
	expectedType := ValMapIntf
	_, err = v.KindSwitch(map[ValKind]KindCase{
		expectedType: func(val interface{}, kind ValKind) (interface{}, error) {
			result, _ = val.(*bwmap.Ordered)
			return nil, nil
		},
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			if len(optDefaultProvider) > 0 {
				fn := optDefaultProvider[0]
				if fn != nil {
					result = fn()
				} else {
					return nil, v.notOfValKindError(ValKindSetFrom(expectedType))
				}
			}
			return nil, nil
		},
	})
	return
}

func (v Holder) MustMapIntf(optDefaultProvider ...func() bwmap.I) (result bwmap.I) {
	var err error
	if result, err = v.MapIntf(optDefaultProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

type ForEachBody func(idx int, key string, hVal Holder) (needBreak bool, err error)

func (v Holder) ForEach(body ForEachBody) (err error) {
	var needBreak bool
	_, err = v.KindSwitch(map[ValKind]KindCase{
		ValMapIntf: func(val interface{}, kind ValKind) (result interface{}, err error) {
			o, _ := val.(bwmap.I)
			for idx, key := range o.Keys() {
				hVal := v.MustKey(key)
				if needBreak, err = body(idx, key, hVal); needBreak || err != nil {
					break
				}
			}
			return
		},
		ValArray: func(val interface{}, kind ValKind) (result interface{}, err error) {
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
	_, _ = v.KindSwitch(map[ValKind]KindCase{
		ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
			m, _ := val.(bwmap.I)
			result = m.Keys(optFilter...)
			return nil, nil
		},
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
	_, _ = v.KindSwitch(map[ValKind]KindCase{
		// ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
		// 	m, _ := val.(bwmap.I)
		// 	_, result = m[key]
		// 	return nil, nil
		// },
		ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
			o, _ := val.(bwmap.I)
			result = o.HasKey(key)
			return nil, nil
		},
	})
	return
}

func (v Holder) DelKey(key string) (result bool) {
	_, _ = v.KindSwitch(map[ValKind]KindCase{
		// ValMap: func(val interface{}, kind ValKind) (interface{}, error) {
		// 	m, _ := val.(map[string]interface{})
		// 	_, result = m[key]
		// 	return nil, nil
		// },
		ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
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
	_, err = v.KindSwitch(map[ValKind]KindCase{
		ValNil: func(val interface{}, kind ValKind) (interface{}, error) {
			var (
				ok  bool
				err error
			)
			if result, ok = defaultVal(optDefaultValProvider); !ok {
				err = v.wrongValError()
			}
			return val, err
		},
		ValMapIntf: func(val interface{}, kind ValKind) (interface{}, error) {
			o, _ := val.(bwmap.I)
			var (
				ok  bool
				err error
			)
			if result, ok = o.Get(key); !ok {
				if result, ok = defaultVal(optDefaultValProvider); !ok {
					// bwdebug.Print("v.Pth", v.Pth)
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
		result = Holder{Val: val, Pth: v.Pth.AppendKey(key)}
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
	_, err = v.KindSwitch(map[ValKind]KindCase{
		// ValMap: func(valMap interface{}, kind ValKind) (interface{}, error) {
		// 	m, _ := valMap.(map[string]interface{})
		// 	m[key] = val
		// 	return nil, nil
		// },
		// ValOrderedMap: func(valOrdered interface{}, kind ValKind) (interface{}, error) {
		// 	o, _ := valOrdered.(*bwmap.Ordered)
		// 	o.Set(key, val)
		// 	return nil, nil
		// },
		ValMapIntf: func(valOrdered interface{}, kind ValKind) (interface{}, error) {
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
		result = Holder{Val: val, Pth: v.Pth.AppendIdx(idx)}
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

func (v Holder) ValidVal(def *Def) (result interface{}, err error) {
	result, err = v.validVal(def, false)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	}
	return
}

func (v Holder) MustValidVal(def *Def) (result interface{}) {
	var err error
	if result, err = v.ValidVal(def); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func (v *Holder) Valid(def *Def) (result Holder, err error) {
	var val interface{}
	val, err = v.validVal(def, false)
	if err != nil {
		err = bwerr.Refine(err, ansi.String("<ansiVal>%s<ansi>::{Error}"), bwjson.Pretty(v.Val))
	} else {
		result = Holder{Val: val, Pth: v.Pth}
	}
	return
}

func (v Holder) MustValid(def *Def) (result Holder) {
	var err error
	if result, err = v.Valid(def); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

type KindCase func(val interface{}, kind ValKind) (interface{}, error)

func (v Holder) KindSwitch(kindCases map[ValKind]KindCase, optDefaultCase ...KindCase) (val interface{}, err error) {
	expects := ValKindSet{}
	for k, _ := range kindCases {
		expects.Add(k)
	}
	val, kind := Kind(v.Val, expects)
	if KindCase, ok := kindCases[kind]; ok {
		val, err = KindCase(val, kind)
	} else if len(optDefaultCase) == 0 {
		vkSet := ValKindSet{}
		for vk := range kindCases {
			vkSet.Add(vk)
		}
		err = v.notOfValKindError(vkSet)
	} else if optDefaultCase[0] != nil {
		val, err = optDefaultCase[0](val, kind)
	}
	return
}

func (v Holder) MustKindSwitch(kindCases map[ValKind]KindCase, optDefaultCase ...KindCase) (val interface{}) {
	var err error
	if val, err = v.KindSwitch(kindCases, optDefaultCase...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

func (v Holder) validVal(def *Def, skipArrayOf bool) (result interface{}, err error) {
	var (
		defKind    ValKind
		val        interface{}
		needReturn bool
	)
	setKind := func(val interface{}, kind ValKind) (interface{}, error) {
		if def.Types.Has(kind) {
			defKind = kind
		}
		return val, nil
	}
	// // bwdebug.Print("v.Val:#v", v.Val)
	if val, err = v.KindSwitch(map[ValKind]KindCase{
		ValNil: func(val interface{}, kind ValKind) (result interface{}, err error) {
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
			if def.Types.Has(ValMap) {
				result = map[string]interface{}{}
				defKind = ValMap
			} else {
				err = v.notOfValKindError(def.Types)
			}
			// // bwdebug.Print("result", result, "defKind", defKind)
			return
		},
		ValBool: setKind,
		ValMap:  setKind,
		ValOrderedMap: func(val interface{}, kind ValKind) (result interface{}, err error) {
			// // bwdebug.Print("!HERE")
			if def.Types.Has(kind) {
				defKind = kind
				result = val
			} else if def.Types.Has(ValMap) {
				defKind = ValMap
				o, _ := val.(*bwmap.Ordered)
				result = o.Map()
			}
			return
		},
		ValString: setKind,
		ValInt: func(val interface{}, kind ValKind) (interface{}, error) {
			if def.Types.Has(ValInt) {
				defKind = ValInt
			} else if def.Types.Has(ValNumber) {
				defKind = ValNumber
			}
			return val, nil
		},
		ValFloat64: func(val interface{}, kind ValKind) (interface{}, error) {
			if def.Types.Has(ValNumber) {
				defKind = ValNumber
			}
			return val, nil
		},
		ValArray: func(val interface{}, kind ValKind) (interface{}, error) {
			if def.Types.Has(ValArray) || !skipArrayOf && def.IsArrayOf {
				defKind = ValArray
			}
			return val, nil
		},
	}, nil); err != nil || needReturn {
		result = val
		return
	}

	if defKind == ValUnknown {
		err = v.notOfValKindError(def.Types)
		return
	}

	switch defKind {
	case ValBool:
	case ValString:
		if def.Enum != nil {
			s, _ := val.(string)
			if !def.Enum.Has(s) {
				err = v.unexpectedEnumValueError(def.Enum)
				return
			}
		}
	case ValInt, ValNumber:
		if !def.Range.Contains(val) {
			err = v.outOfRangeError(def.Range)
			return
		}

	case ValMap, ValOrderedMap:
		var m map[string]interface{}
		if defKind == ValMap {
			m, _ = val.(map[string]interface{})
		} else {
			o, _ := val.(*bwmap.Ordered)
			m = o.Map()
		}
		if def.KeysDef != nil {
			unexpectedKeys := bwmap.MustUnexpectedKeys(val, def.KeysDef)
			for key, keyDef := range def.KeysDef {
				// bwdebug.Print("v:json", v, "key", key, "m:json", m, "keyDef:json", keyDef)
				if err = mapHelper(v.Pth, m, key, &keyDef); err != nil {
					return
				}
			}
			if unexpectedKeys != nil {
				if def.ElemDef == nil {
					err = v.unexpectedKeysError(unexpectedKeys)
					return
				} else {
					for _, key := range unexpectedKeys.ToSlice() {
						if err = mapHelper(v.Pth, m, key, def.ElemDef); err != nil {
							return
						}
					}
				}
			}
			var requiredKeys []string
			for key, keyDef := range def.KeysDef {
				if !keyDef.IsOptional {
					requiredKeys = append(requiredKeys, key)
				}
			}
			if len(requiredKeys) > 0 {
				for _, key := range requiredKeys {
					if _, ok := m[key]; !ok {
						err = v.hasNoKeyError(key)
						return
					}
				}
			}
		}
		if def.ElemDef != nil {
			for k := range m {
				if def.KeysDef != nil {
					if _, ok := def.KeysDef[k]; ok {
						continue
					}
				}
				if err = mapHelper(v.Pth, m, k, def.ElemDef); err != nil {
					return
				}
			}
		}
	case ValArray:
		if !skipArrayOf && def.IsArrayOf {
			v.Val = val
			if val, err = v.arrayHelper(def, true); err != nil {
				return
			}
		} else {
			elemDef := def.ArrayElemDef
			if elemDef == nil {
				elemDef = def.ElemDef
			}
			if elemDef != nil {
				v.Val = val
				if val, err = v.arrayHelper(elemDef, false); err != nil {
					return
				}
			}
		}
	}

	if !skipArrayOf && def.IsArrayOf && defKind != ValArray {
		val = []interface{}{val}
	}

	result = val
	return
}

func mapHelper(path bw.ValPath, m map[string]interface{}, key string, elemDef *Def) (err error) {
	vp, _ := Holder{Val: m, Pth: path}.Key(key)
	// bwdebug.Print("vp:json", vp, "path", path, "key", key)
	var val interface{}
	if val, err = vp.validVal(elemDef, false); err != nil {
		return
	} else if val != nil {
		m[key] = val
		// // bwdebug.Print("m", m)
	}
	return
}

func (v Holder) arrayHelper(elemDef *Def, skipArrayOf bool) (result interface{}, err error) {
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
			if _, err = h.KindSwitch(map[ValKind]KindCase{
				ValString: func(val interface{}, kind ValKind) (interface{}, error) {
					s, _ := val.(string)
					result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Name: s})
					return val, nil
				},
				ValInt: func(val interface{}, kind ValKind) (interface{}, error) {
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
	_, err = v.KindSwitch(map[ValKind]KindCase{
		ValArray: func(val interface{}, kind ValKind) (interface{}, error) {
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
