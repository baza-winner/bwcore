package bwval

import (
	"fmt"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwstr"

	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwset"
)

// ============================================================================

var (
	ansiCanNotBeMixed                       string
	ansiUnexpectedKey                       string
	ansiNoRequiredKey                       string
	fmtKeyIsSpecifiedSoDefTypeExpectsToHave string
	fmtAnsiVal                              string
	fmtKeyMustBeSpecifiedFirst              string
	fmtKeyCanNotbeSpecified                 string
	fmtKeyMustBeLastSpecified               string
	fmtKeyMustBeSpecified                   string
	ansiArrayOfMustBeFollowed               string
)

const ArrayOf = "ArrayOf"

func init() {
	ansiCanNotBeMixed = ansi.String("<ansiErr>%s<ansi> can not be mixed with <ansiVal>%s")
	ansiUnexpectedKey = ansi.String("unexpected key `<ansiErr>%s<ansi>`")
	ansiNoRequiredKey = ansi.String("no <ansiErr>required<ansi> key `<ansiErr>%s<ansi>`")
	fmtKeyIsSpecifiedSoDefTypeExpectsToHave = "key <ansiVar>%d<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have "
	fmtAnsiVal = "<ansiVal>%s<ansi>"
	fmtKeyMustBeSpecifiedFirst = "key <ansiVar>%s<ansi> must be specified first"
	fmtKeyCanNotbeSpecified = "while <ansiVar>%s <ansiVal>%v<ansi>, key <ansiVar>%s<ansi> can not be specified"
	fmtKeyMustBeLastSpecified = "key <ansiVar>%s<ansi> must be LAST specified, but found key <ansiErr>%s<ansi> after it"
	fmtKeyMustBeSpecified = "key <ansiVar>%s<ansi> must be specified"
	ansiArrayOfMustBeFollowed = "<ansiType>" + ArrayOf + "<ansi> must be followed by another <ansiVar>Type<ansi>"
}

var defSupportedTypes ValKindSet

func init() {
	defSupportedTypes = ValKindSetFrom(
		ValString,
		ValBool,
		ValInt,
		ValNumber,
		ValMap,
		ValArray,
	)
}

const (
	defType         = "type"
	defIsOptional   = "isOptional"
	defDefault      = "default"
	defEnum         = "enum"
	defRange        = "range"
	defKeysDef      = "keysDef"
	defElemDef      = "elemDef"
	defArrayElemDef = "arrayElemDef"
)

func ParseDefTemplate(p bwparse.I, optBaseProvider ...bw.ValPathProvider) (result *DefTemplate, status Status) {
	result = &DefTemplate{}
	defer func() {
		if !status.IsOK() {
			result = nil
		}
	}()

	var base bw.ValPath
	if len(optBaseProvider) > 0 {
		base = MustPath(optBaseProvider[0])
	} else {
		base = bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemVar, Name: "Def"}}
	}
	string2type := map[string]ValKind{}
	for _, t := range defSupportedTypes.ToSlice() {
		s := t.String()
		string2type[s] = t
	}
	var (
		tp                  ValKind
		ok                  bool
		val                 interface{}
		isOptionalSpecified bool
	)
	result.Types = ValKindSet{}
	addType := func(on On, tp ValKind) (err error) {
		switch tp {
		case ValNumber, ValInt:
			dontMix := ValKindSetFrom(ValInt, ValNumber)
			for _, vk := range dontMix.ToSlice() {
				if result.Types.Has(vk) {
					err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, tp, vk)})
					return
				}
			}
		case ValArray:
			if result.IsArrayOf {
				err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, ValArray, ArrayOf)})
				return
			}
		}
		result.Types.Add(tp)
		return
	}
	processType := func(on On, s string) (ok bool, err error) {
		if tp, ok = string2type[s]; ok {
			err = addType(on, tp)
		} else if ok = s == ArrayOf; ok {
			if result.Types.Has(ValArray) {
				err = p.Error(bwparse.E{Start: on.Start, Fmt: bw.Fmt(ansiCanNotBeMixed, ArrayOf, ValArray)})
			} else {
				result.IsArrayOf = true
			}
		}
		return
	}
	// validKeys := bwset.StringFrom("type", "range", "enum", "keysDef", "elemDef", "arrayElemDef", "default", "isOptional")
	validKeys := bwset.StringFrom(defType, defRange, defEnum, defKeysDef, defElemDef, defArrayElemDef, defDefault, defIsOptional)
	setForType := func(opt Opt) Opt {
		opt.KindSet.Add(ValId, ValString, ValArray)
		opt.OnId = func(on On, s string) (val interface{}, ok bool, err error) {
			ok, err = processType(on, s)
			return
		}
		opt.OnValidateString = func(on On, s string) (err error) {
			if ok, err = processType(on, s); !ok && err == nil {
				err = bwparse.Unexpected(p, on.Start)
			}
			return
		}
		opt.OnParseArrayElem = func(on On, vals []interface{}) (outVals []interface{}, status Status) {
			var val interface{}
			kindSet := on.Opt.KindSet
			defer func() { on.Opt.KindSet = kindSet }()
			on.Opt.KindSet = ValKindSetFrom(ValId, ValString)
			if val, status = ParseVal(p, *on.Opt); status.IsOK() {
				outVals = append(vals, val)
			}
			return
		}
		opt.OnValidateArray = func(on On, vals []interface{}) (err error) {
			if len(vals) == 0 {
				// err = bwparse.Expects(p, bwparse.Unexpected(p, on.Start), "non empty <ansiType>Array<ansi>")
				err = bwparse.Expects(p, bwparse.Unexpected(p, on.Start), "non empty "+ValArray.AnsiString())
			}
			return
		}
		return opt
	}
	checkArrayOf := func(start *bwparse.Start) (err error) {
		if result.IsArrayOf && len(result.Types) == 0 {
			err = p.Error(bwparse.E{
				Start: start,
				// Fmt:   bw.A{Fmt: "<ansiType>ArrayOf<ansi> must be followed by another <ansiVar>Type<ansi>"},
				Fmt: bw.A{Fmt: ansiArrayOfMustBeFollowed},
			})
		}
		return err
	}
	if val, status = ParseVal(p, setForType(Opt{
		ValFalse: true,
		KindSet:  ValKindSetFrom(ValMap),
		OnValidateMap: func(on On, m bwmap.I) (err error) {
			if len(result.Types) == 0 {
				err = p.Error(bwparse.E{
					Start: on.Start,
					// Fmt:   bw.A{Fmt: "key <ansiVar>type<ansi> must be specified"},
					Fmt: bw.Fmt(fmtKeyMustBeSpecified, defType),
				})
			}
			return
		},
		OnValidateMapKey: func(on On, m bwmap.I, key string) (err error) {
			if !validKeys.Has(key) {
				err = p.Error(bwparse.E{
					Start: on.Start,
					Fmt:   bw.Fmt(ansiUnexpectedKey, on.Start.Suffix()),
				})
			} else {
				switch key {
				case "default":
					if len(result.Types) == 0 {
						err = p.Error(bwparse.E{
							Start: on.Start,
							// Fmt:   bw.A{Fmt: "key <ansiVar>type<ansi> must be specified first"},
							Fmt: bw.Fmt(fmtKeyMustBeSpecifiedFirst, defType),
						})
					} else if isOptionalSpecified && !result.IsOptional {
						err = p.Error(bwparse.E{
							Start: on.Start,
							Fmt:   bw.Fmt(fmtKeyCanNotbeSpecified, defIsOptional, false, key),
						})
					}
				default:
					if result.Default != nil {
						err = p.Error(bwparse.E{
							Start: on.Start,
							Fmt:   bw.Fmt(fmtKeyMustBeLastSpecified, defDefault, key),
						})
					}
				}
			}
			return
		},
		OnParseMapElem: func(on On, m bwmap.I, key string) (status Status) {
			switch key {
			case defType:
				if _, status = ParseVal(p, setForType(Opt{KindSet: ValKindSet{}})); status.IsOK() {
					validKeys = bwset.StringFrom(defIsOptional, defDefault)
					for _, vk := range result.Types.ToSlice() {
						switch vk {
						case ValString:
							validKeys.Add(defEnum)
						case ValInt, ValNumber:
							validKeys.Add(defRange)
						case ValMap, ValOrderedMap:
							validKeys.Add(defKeysDef)
							validKeys.Add(defElemDef)
						case ValArray:
							validKeys.Add(defArrayElemDef)
							validKeys.Add(defElemDef)
						}
					}
					if status.Err = checkArrayOf(status.Start); status.Err == nil {
						err := func(key string, kinds ...ValKind) error {
							return p.Error(bwparse.E{
								Start: status.Start,
								Fmt: bw.Fmt(
									fmtKeyIsSpecifiedSoDefTypeExpectsToHave+
										bwstr.SmartJoin(bwstr.A{
											Source: bwstr.SS{
												SS: ValKinds(kinds).Strings(),
												Preformat: func(s string) string {
													return fmt.Sprintf(fmtAnsiVal, s)
												},
											},
										}),
									key, defType,
								),
							})
						}
						if result.Enum != nil && !result.Types.Has(ValString) {
							status.Err = err(defEnum, ValString)
							// status.Err = p.Error(bwparse.E{
							// 	Start: status.Start,
							// 	Fmt: bw.A{
							// 		Fmt:  "key <ansiVar>%d<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have <ansiVal>%s<ansi>",
							// 		Args: []interface{}{defEnum, defType, ValString.String()},
							// 	},
							// })

						} else if result.Range != nil && !(result.Types.Has(ValInt) || result.Types.Has(ValNumber)) {
							status.Err = err(defRange, ValInt, ValNumber)
							// status.Err = p.Error(bwparse.E{
							// 	Start: status.Start,
							// 	Fmt: bw.A{
							// 		Fmt:  "key <ansiVar>%s<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have <ansiVal>%s<ansi> or <ansiVal>%s<ansi>",
							// 		Args: []interface{}{defRange, defType, ValInt.String(), ValNumber.String()},
							// 	},
							// })
						} else if result.KeysDef != nil && !result.Types.Has(ValMap) {
							status.Err = err(defKeysDef, ValMap)
							// status.Err = p.Error(bwparse.E{
							// 	Start: status.Start,
							// 	Fmt: bw.A{
							// 		Fmt:  "key <ansiVar>%s<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have <ansiVal>%s<ansi>",
							// 		Args: []interface{}{defKeysDef, defType, ValMap.String()},
							// 	},
							// })
						} else if result.ArrayElemDef != nil && !result.Types.Has(ValArray) {
							status.Err = err(defArrayElemDef, ValArray)
							// status.Err = p.Error(bwparse.E{
							// 	Start: status.Start,
							// 	Fmt: bw.A{
							// 		Fmt:  "key <ansiVar>%s<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have <ansiVal>%s<ansi>",
							// 		Args: []interface{}{defArrayElemDef, defType, ValArray.String()},
							// 	},
							// })
						} else if result.ElemDef != nil &&
							!(result.Types.Has(ValMap) ||
								result.ArrayElemDef == nil && result.Types.Has(ValArray)) {
							kinds := []ValKind{ValMap}
							if result.ArrayElemDef == nil {
								kinds = append(kinds, ValArray)
							}
							status.Err = err(defElemDef, kinds...)
							// var suffix string
							// if result.ArrayElemDef == nil {
							// 	suffix = " or <ansiVal>Array<ansi>"
							// }
							// status.Err = p.Error(bwparse.E{
							// 	Start: status.Start,
							// 	Fmt: bw.A{
							// 		Fmt:  "key <ansiVar>%s<ansi> is specified, so value of key <ansiVar>%s<ansi> expects to have <ansiVal>%s<ansi>" + suffix,
							// 		Args: []interface{}{defElemDef, defType, ValMap.String()},
							// 	},
							// })
						}
					}
				}
			case defEnum:
				onValidateString := func(on On, s string) (err error) {
					result.Enum = append(result.Enum, EnumTemplateItem{S: s})
					return
				}
				onValidatePath := func(on On, path bw.ValPathTemplate) (err error) {
					result.Enum = append(result.Enum, EnumTemplateItem{IsPath: true, Path: path})
					return
				}
				_, status = ParseVal(p, Opt{
					KindSet:          ValKindSetFrom(ValString, ValArray, ValPath),
					OnValidatePath:   onValidatePath,
					OnValidateString: onValidateString,
					OnParseArrayElem: func(on On, vals []interface{}) (outVals []interface{}, status Status) {
						var val interface{}
						if val, status = ParseVal(p, Opt{
							KindSet:          ValKindSetFrom(ValString, ValPath),
							OnValidatePath:   onValidatePath,
							OnValidateString: onValidateString,
						}); status.IsOK() {
							outVals = append(vals, val)
							if s, ok := val.(string); ok {
								err = onValidateString(on, s)
								// result.Enum = append(result.Enum, EnumTemplateItem{S: s})
							} else {
								path, _ := val.(bw.ValPathTemplate)
								err = onValidatePath(on, path)
							}
						}
						// var s string
						// if s, status = ParseString(p, *on.Opt); status.IsOK() {
						// 	// if result.Enum == nil {
						// 	// 	result.Enum = []EnumTemplateItem{}
						// 	// }
						// 	// result.Enum.Add(s)
						// 	result.Enum = append(result.Enum, EnumTemplateItem{S: s})
						// 	outVals = append(vals, s)
						// } else if status.Err == nil {
						// 	// status.Err = bwparse.Expects(p, nil, "<ansiType>String<ansi>")
						// 	status.Err = bwparse.Expects(p, nil, ValString.AnsiString())
						// }
						return
					},
				})
			case defRange:
				_, status = ParseVal(p, Opt{
					KindSet: ValKindSetFrom(ValRange, ValPath),
					OnValidatePath: func(on On, path bw.ValPathTemplate) (err error) {
						result.Range = DefTemplateRange{IsPath: true, Path: path}
						// result.Enum = append(result.Enum, EnumTemplateItem{IsPath: true, Path: path})
						return
					},
					OnValidateRange: func(on On, rng *RangeTemplate) (err error) {
						result.Range = DefTemplateRange{RangeTemplate: rng}
						// result.Enum = append(result.Enum, EnumTemplateItem{IsPath: true, Path: path})
						return
					},
				})
				// if result.Range, status = ParseRange(p, Opt{KindSet: ValKindSetFrom(ValNumber)}); !status.OK && status.Err == nil {
				// 	// status.Err = bwparse.Expects(p, nil, "<ansiType>Range<ansi>")
				// 	status.Err = bwparse.Expects(p, nil, ValRange.AnsiString())
				// }
			case defKeysDef:
				if _, status = ParseMap(p, Opt{
					OnParseMapElem: func(on On, m bwmap.I, key string) (status Status) {
						if result.KeysDef == nil {
							result.KeysDef = map[string]Def{}
						}
						var def *Def
						if def, status = ParseDef(p, base.AppendKey(key)); status.IsOK() {
							result.KeysDef[key] = *def
							m.Set(key, nil)
						} else if status.Err == nil {
							// status.Err = bwparse.Expects(p, nil, "<ansiType>Def<ansi>")
							status.Err = bwparse.Expects(p, nil, ValDef.AnsiString())
						}
						return
					},
				}); !status.OK && status.Err == nil {
					// status.Err = bwparse.Expects(p, nil, "<ansiType>Map<ansi>")
					status.Err = bwparse.Expects(p, nil, ValMap.AnsiString())
				}
			case defElemDef, defArrayElemDef:
				var def *Def
				if def, status = ParseDef(p, base.AppendKey(key)); status.IsOK() {
					if key == defElemDef {
						result.ElemDef = def
					} else {
						result.ArrayElemDef = def
					}
				} else if status.Err == nil {
					// status.Err = bwparse.Expects(p, nil, "<ansiType>Def<ansi>")
					status.Err = bwparse.Expects(p, nil, ValDef.AnsiString())
				}
			case defIsOptional:
				result.IsOptional, status = ParseBool(p)
				isOptionalSpecified = true
			case defDefault:
				result.Default, status = ParseValByDef(p, result, base.AppendKey(key))
				if result.Default != nil {
					result.IsOptional = true
				}
			}
			m.Set(key, nil)
			return
		},
	})); status.IsOK() {
		if _, ok := val.(map[string]interface{}); !ok {
			status.Err = checkArrayOf(status.Start)
		}
	}
	return
}

// ============================================================================

func ParseValByDef(p bwparse.I, def *Def, optBase ...bw.ValPath) (result interface{}, status Status) {
	var base bw.ValPath
	if len(optBase) > 0 {
		base = optBase[0]
	}
	return parseValByDef(p, def, base, false)
}

func parseValByDef(p bwparse.I, def *Def, base bw.ValPath, skipArrayOf bool) (result interface{}, status Status) {

	opt := Opt{KindSet: ValKindSet{}}
	opt.KindSet.AddSet(def.Types)

	if opt.KindSet.Has(ValMap) || opt.KindSet.Has(ValOrderedMap) {
		opt.OnValidateMapKey = func(on On, m bwmap.I, key string) (err error) {
			if def.ElemDef == nil && def.KeysDef != nil {
				if _, ok := def.KeysDef[key]; !ok {
					err = p.Error(bwparse.E{
						Start: on.Start,
						// Fmt:   bw.Fmt(ansi.String("unexpected key `<ansiErr>%s<ansi>`"), on.Start.Suffix()),
						Fmt: bw.Fmt(ansiUnexpectedKey, on.Start.Suffix()),
					})
				}
			}
			return
		}
		opt.OnParseMapElem = func(on On, m bwmap.I, key string) (status Status) {
			var keyDef *Def
			if def.KeysDef != nil {
				if v, ok := def.KeysDef[key]; ok {
					keyDef = &v
				}
			}
			if keyDef == nil {
				keyDef = def.ElemDef
			}
			if keyDef == nil {
				keyDef = &Def{Types: defSupportedTypes}
			}
			var val interface{}
			if val, status = parseValByDef(p, keyDef, base.AppendKey(key), skipArrayOf); status.IsOK() {
				m.Set(key, val)
			}
			return
		}
		opt.OnValidateMap = func(on On, m bwmap.I) (err error) {
			if def.KeysDef != nil {
				var requiredKeys []string
				for key, keyDef := range def.KeysDef {
					if !keyDef.IsOptional {
						requiredKeys = append(requiredKeys, key)
					}
				}
				if len(requiredKeys) > 0 {
					for _, key := range requiredKeys {
						if !m.HasKey(key) {
							err = p.Error(bwparse.E{
								Start: on.Start,
								// Fmt:   bw.Fmt(ansi.String("no <ansiErr>required<ansi> key `<ansiErr>%s<ansi>`"), key),
								Fmt: bw.Fmt(ansiNoRequiredKey, key),
							})
							return
						}
					}
				}
			}
			return
		}
	}

	if hasArray := opt.KindSet.Has(ValArray); hasArray || !skipArrayOf && def.IsArrayOf {
		elemDef := func() *Def {
			if !hasArray {
				return def
			} else {
				var elemDef *Def
				if def.ArrayElemDef != nil {
					elemDef = def.ArrayElemDef
				}
				if elemDef == nil {
					elemDef = def.ElemDef
				}
				if elemDef == nil {
					elemDef = &Def{Types: defSupportedTypes}
				}
				return elemDef
			}
		}
		subSkipArrayOf := skipArrayOf
		if !hasArray {
			opt.KindSet.Add(ValArray)
			subSkipArrayOf = true
		}
		opt.OnParseArrayElem = func(on On, vals []interface{}) (outVals []interface{}, status Status) {
			var val interface{}
			if val, status = parseValByDef(p, elemDef(), base.AppendIdx(len(vals)), subSkipArrayOf); status.IsOK() {
				outVals = append(vals, val)
			}
			return
		}
		opt.OnValidateArrayOfStringElem = func(on On, ss []string, s string) (err error) {
			h := Holder{Val: s, Pth: base.AppendIdx(len(ss))}
			if _, err = h.validVal(elemDef(), subSkipArrayOf); err != nil {
				// bwdebug.Print("err", err)
				err = p.Error(bwparse.E{
					Start: on.Start,
					Fmt:   bwerr.Err(err),
				})
			}
			return
		}
	}

	if result, status = ParseVal(p, opt); status.IsOK() {
		h := Holder{Val: result, Pth: base}
		if result, status.Err = h.validVal(def, skipArrayOf); status.Err != nil {
			status.Err = p.Error(bwparse.E{
				Start: status.Start,
				Fmt:   bwerr.Err(status.Err),
			})
		}
	}
	return
}

// ============================================================================
