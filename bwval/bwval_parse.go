package bwval

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwmap"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwset"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

// type PosInfo struct {
//   isEOF       bool
//   rune        rune
//   pos         int
//   line        uint
//   col         uint
//   prefix      string
//   prefixStart int
//   justParsed  interface{}
// }

// func (p PosInfo) IsEOF() bool {
//   return p.isEOF
// }

// func (p PosInfo) Rune() rune {
//   return p.rune
// }

// // ============================================================================

// func (start Start) Suffix() string {
//   return start.suffix
// }

// // ============================================================================

// type bwparse.I interface {
//   FileSpec() string
//   Close() error
//   Curr() *PosInfo
//   Forward(count uint)
//   Error(a E) error
//   LookAhead(ofs uint) *PosInfo
//   Start() *bwparse.Start
//   Stop(start *bwparse.Start)
// }

// ============================================================================

type On struct {
	P     bwparse.I
	Start *bwparse.Start
	Opt   *Opt
}

type NonNegativeNumberFunc func(rangeLimitKind RangeLimitKind) bool
type IdFunc func(on On, s string) (val interface{}, ok bool, err error)

type ValidateMapKeyFunc func(on On, m bwmap.I, key string) (err error)
type ParseMapElemFunc func(on On, m bwmap.I, key string) (status Status)
type ValidateMapFunc func(on On, m bwmap.I) (err error)

type ParseArrayElemFunc func(on On, vals []interface{}) (outVals []interface{}, status Status)
type ValidateArrayFunc func(on On, vals []interface{}) (err error)

type ValidateNumberFunc func(on On, n *bwtype.Number) (err error)
type ValidateRangeFunc func(on On, rng *bwtype.Range) (err error)
type ValidatePathFunc func(on On, path bw.ValPath) (err error)

type ValidateStringFunc func(on On, s string) (err error)

type ValidateArrayOfStringElemFunc func(on On, ss []string, s string) (err error)

// ============================================================================

type RangeLimitKind uint8

const (
	RangeLimitNone RangeLimitKind = iota
	RangeLimitMin
	RangeLimitMax
)

type Opt struct {
	ExcludeKinds bool
	KindSet      bwtype.ValKindSet

	RangeLimitMinOwnKindSet bwtype.ValKindSet // empty means 'inherits'
	RangeLimitMaxOwnKindSet bwtype.ValKindSet // empty means 'inherits'

	Path bw.ValPath

	StrictId          bool
	IdVals            map[string]interface{}
	OnId              IdFunc
	NonNegativeNumber NonNegativeNumberFunc

	IdNil   bwset.String
	IdFalse bwset.String
	IdTrue  bwset.String

	OnValidateMapKey ValidateMapKeyFunc
	OnParseMapElem   ParseMapElemFunc
	OnValidateMap    ValidateMapFunc

	OnParseArrayElem ParseArrayElemFunc
	OnValidateArray  ValidateArrayFunc

	OnValidateString            ValidateStringFunc
	OnValidateArrayOfStringElem ValidateArrayOfStringElemFunc

	OnValidateNumber ValidateNumberFunc
	OnValidateRange  ValidateRangeFunc
	OnValidatePath   ValidatePathFunc

	ValFalse bool
}

// ============================================================================

// type P struct {
//   pp            bwrune.ProviderProvider
//   prov          bwrune.Provider
//   curr          *PosInfo
//   next          []*PosInfo
//   preLineCount  uint
//   postLineCount uint
//   starts        map[int]*bwparse.Start
// }

// func ParseFrom(pp bwrune.ProviderProvider, opt ...map[string]interface{}) (result *P, err error) {
//   var p bwrune.Provider
//   if p, err = pp.Provider(); err != nil {
//     return
//   }
//   defer func() {
//     if err != nil {
//       p.Close()
//     }
//   }()
//   result = &P{
//     prov:          p,
//     curr:          &PosInfo{pos: -1, line: 1},
//     next:          []*PosInfo{},
//     preLineCount:  3,
//     postLineCount: 3,
//   }
//   if len(opt) > 0 {
//     m := opt[0]
//     if m != nil {
//       keys := bwset.String{}
//       if i, ok := optKeyUint(m, "preLineCount", &keys); ok {
//         result.preLineCount = i
//       }
//       if i, ok := optKeyUint(m, "postLineCount", &keys); ok {
//         result.postLineCount = i
//       }
//       if unexpectedKeys := bwmap.MustUnexpectedKeys(m, keys); len(unexpectedKeys) > 0 {
//         err = bwerr.From(ansiOptHasUnexpectedKeys, bwjson.Pretty(m), unexpectedKeys)
//       }
//     }
//   }
//   return
// }

// func MustParseFrom(pp bwrune.ProviderProvider, opt ...map[string]interface{}) (result *P) {
//   var err error
//   if result, err = ParseFrom(pp, opt...); err != nil {
//     bwerr.PanicErr(err)
//   }
//   return
// }

// const bwparse.Initial uint = 0

// func (p *P) FileSpec() string {
//   return p.prov.FileSpec()
// }

// func (p *P) Close() error {
//   return p.prov.Close()
// }

// func (p *P) Curr() *PosInfo {
//   return p.curr
// }

// func (p *P) Forward(count uint) {
//   if p.curr.pos < 0 || count > 0 && !p.curr.isEOF {
//     if count <= 1 {
//       p.forward()
//     } else {
//       for ; count > 0; count-- {
//         p.forward()
//       }
//     }
//   }
// }

// func (p *P) LookAhead(ofs uint) (result *PosInfo) {
//   result = p.curr
//   if ofs > 0 {
//     idx := len(p.next) - int(ofs)
//     if idx >= 0 {
//       result = p.next[idx]
//     } else {
//       var ps PosInfo
//       if len(p.next) > 0 {
//         ps = *p.next[0]
//       } else {
//         ps = *p.curr
//       }
//       var lookahead []PosInfo
//       for i := idx; i < 0 && !ps.isEOF; i++ {
//         ps = p.pullRune(ps)
//         lookahead = append(lookahead, ps)
//       }
//       var newNext []*PosInfo
//       for i := len(lookahead) - 1; i >= 0; i-- {
//         newNext = append(newNext, &lookahead[i])
//       }
//       p.next = append(newNext, p.next...)
//       if len(p.next) > 0 {
//         result = p.next[0]
//       }
//     }
//   }
//   return
// }

// func (p *P) Start() (result *bwparse.Start) {
//   p.Forward(bwparse.Initial)
//   var ok bool
//   if result, ok = p.starts[p.curr.pos]; !ok {
//     result = &Start{ps: p.curr}
//     if p.starts == nil {
//       p.starts = map[int]*bwparse.Start{}
//     }
//     p.starts[p.curr.pos] = result
//   }
//   return
// }

// type E struct {
//   Start *bwparse.Start
//   Fmt   bw.bwparse.I
// }

// func (p *P) Error(a E) error {
//   var start Start
//   if a.Start == nil {
//     start = Start{ps: p.curr}
//   } else {
//     start = *a.Start
//   }
//   var msg string
//   if start.ps.pos < p.curr.pos {
//     if a.Fmt != nil {
//       msg = bw.Spew.Sprintf(a.Fmt.FmtString(), a.Fmt.FmtArgs()...)
//     } else {
//       msg = fmt.Sprintf(ansiUnexpectedWord, start.suffix)
//     }
//   } else if !p.curr.isEOF {
//     msg = fmt.Sprintf(ansiUnexpectedChar, p.curr.rune, p.curr.rune)
//   } else {
//     msg = ansiUnexpectedEOF
//   }
//   return bwerr.From(msg + suffix(&bwparse.Proxy{p: p}, start, p.postLineCount))
// }

// // ============================================================================

// func Unexpected(p bwparse.I, optStart ...*bwparse.Start) error {
//   var a E
//   if len(optStart) > 0 {
//     a.Start = optStart[0]
//   }
//   return p.Error(a)
// }

// func bwparse.Expects(p bwparse.I, err error, what string) error {
//   if err == nil {
//     err = Unexpected(p)
//   }
//   return bwerr.Refine(err, "expects %s instead of {Error}", what)
// }

// func bwparse.ExpectsSpace(p bwparse.I) error {
//   return bwparse.Expects(p, nil, ansi.String(ansiVarSpace))
// }

// // ============================================================================

// func bwparse.CheckNotEOF(p bwparse.I) (err error) {
//   if p.Curr().isEOF {
//     err = Unexpected(p)
//   }
//   return
// }

// // ============================================================================

// func bwparse.CanSkipRunes(p bwparse.I, rr ...rune) bool {
//   for i, r := range rr {
//     if pi := p.LookAhead(uint(i)); pi.isEOF || pi.rune != r {
//       return false
//     }
//   }
//   return true
// }

// func bwparse.SkipRunes(p bwparse.I, rr ...rune) (ok bool) {
//   if ok = bwparse.CanSkipRunes(p, rr...); ok {
//     p.Forward(uint(len(rr)))
//   }
//   return
// }

// // ============================================================================

// func IsDigit(r rune) bool {
//   return '0' <= r && r <= '9'
// }

// func IsLetter(r rune) bool {
//   return r == '_' || unicode.IsLetter(r)
// }

// func IsPunctOrSymbol(r rune) bool {
//   return unicode.IsPunct(r) || unicode.IsSymbol(r)
// }

// // ============================================================================

// const (
//   bwparse.TillNonEOF bool = false
//   bwparse.TillEOF    bool = true
// )

// func bwparse.SkipSpace(p bwparse.I, tillEOF bool) (ok bool, err error) {
//   p.Forward(bwparse.Initial)
// REDO:
//   for !p.Curr().isEOF && unicode.IsSpace(p.Curr().rune) {
//     ok = true
//     p.Forward(1)
//   }
//   if p.Curr().isEOF && !tillEOF {
//     err = Unexpected(p)
//     return
//   }
//   if bwparse.CanSkipRunes(p, '/', '/') {
//     ok = true
//     p.Forward(2)
//     for !p.Curr().isEOF && p.Curr().rune != '\n' {
//       p.Forward(1)
//     }
//     if !p.Curr().isEOF {
//       p.Forward(1)
//     }
//     goto REDO
//   } else if bwparse.CanSkipRunes(p, '/', '*') {
//     ok = true
//     p.Forward(2)
//     for !p.Curr().isEOF && !bwparse.CanSkipRunes(p, '*', '/') {
//       p.Forward(1)
//     }
//     if !p.Curr().isEOF {
//       p.Forward(2)
//     }
//     goto REDO
//   }
//   if tillEOF && !p.Curr().isEOF {
//     err = Unexpected(p)
//   }
//   return
// }

// ============================================================================

type Status struct {
	Start *bwparse.Start
	OK    bool
	Err   error
}

func (v Status) IsOK() bool {
	return v.OK && v.Err == nil
}

func (v Status) NoErr() bool {
	return v.Err == nil
}

func (v *Status) UnexpectedIfErr(p bwparse.I) {
	if v.Err != nil {
		v.Err = p.Error(bwparse.E{v.Start, bwerr.Err(v.Err)})
	}
}

// ============================================================================

func ParseNil(p bwparse.I, optOpt ...Opt) (status Status) {
	opt := getOpt(optOpt)

	ss := []string{"nil"}
	if len(opt.IdNil) > 0 {
		ss = append(ss, opt.IdNil.ToSliceOfStrings()...)
	}

	var needForward uint
	if needForward, status.OK = isOneOfId(p, ss); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(needForward)
	}
	return
}

// ============================================================================

func ParseBool(p bwparse.I, optOpt ...Opt) (result bool, status Status) {
	opt := getOpt(optOpt)

	ss := []string{"true"}
	if len(opt.IdTrue) > 0 {
		ss = append(ss, opt.IdTrue.ToSliceOfStrings()...)
	}

	var needForward uint
	if needForward, status.OK = isOneOfId(p, ss); status.OK {
		result = true
	} else {
		ss = []string{"false"}
		if len(opt.IdFalse) > 0 {
			ss = append(ss, opt.IdFalse.ToSliceOfStrings()...)
		}
		needForward, status.OK = isOneOfId(p, ss)
	}
	if status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(needForward)
	}
	return
}

// ============================================================================

func ParseId(p bwparse.I, optOpt ...Opt) (result string, status Status) {
	opt := getOpt(optOpt)
	r := p.Curr().Rune()
	var isId func(r rune) bool
	if opt.StrictId {
		isId = func(r rune) bool { return bw.IsLetter(r) }
	} else {
		isId = func(r rune) bool { return bw.IsLetter(r) || bw.IsDigit(r) || r == '-' || r == '/' || r == '.' }
	}
	if status.OK = isId(r); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		var while func(r rune) bool
		if opt.StrictId {
			while = func(r rune) bool { return bw.IsLetter(r) || bw.IsDigit(r) }
		} else {
			while = func(r rune) bool { return bw.IsLetter(r) || bw.IsDigit(r) || r == '-' || r == '/' || r == '.' }
		}
		for !p.Curr().IsEOF() && while(r) {
			result += string(r)
			p.Forward(1)
			r = p.Curr().Rune()
		}
	}
	return
}

// func isReservedRune(r rune) bool {
//  return r == ',' ||
//    r == '{' || r == '<' || r == '[' || r == '(' ||
//    r == '}' || r == '>' || r == ']' || r == ')' ||
//    r == ':' || r == '=' ||
//    r == '"' || r == '\''
// }

// ============================================================================

func ParseString(p bwparse.I, optOpt ...Opt) (result string, status Status) {
	delimiter := p.Curr().Rune()
	if status.OK = bwparse.CanSkipRunes(p, '"') || bwparse.CanSkipRunes(p, '\''); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		expectEscapedContent := false
		b := true
		for status.NoErr() {
			r := p.Curr().Rune()
			if !expectEscapedContent {
				if p.Curr().IsEOF() {
					b = false
				} else if bwparse.SkipRunes(p, delimiter) {
					break
				} else if bwparse.SkipRunes(p, '\\') {
					expectEscapedContent = true
					continue
				}
			} else if !(r == '"' || r == '\'' || r == '\\') {
				r, b = EscapeRunes[r]
				b = b && delimiter == '"'
			}
			if !b {
				status.Err = bwparse.Unexpected(p)
			} else {
				result += string(r)
				p.Forward(1)
			}
			expectEscapedContent = false
		}
	}

	return
}

var EscapeRunes = map[rune]rune{
	'a': '\a',
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
	'v': '\v',
}

// ============================================================================

func ParseInt(p bwparse.I, optOpt ...Opt) (result int, status Status) {
	opt := getOpt(optOpt)
	result, status = parseInt(p, opt, RangeLimitNone)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

// ============================================================================

func ParseUint(p bwparse.I, optOpt ...Opt) (result uint, status Status) {
	var justParsed numberResult
	curr := p.Curr()
	if justParsed, status.OK = curr.JustParsed.(numberResult); status.OK {
		if result, status.OK = bwtype.Uint(justParsed.n.Val()); status.OK {
			status.Start = justParsed.start
			p.Forward(uint(len(justParsed.start.Suffix())))
			return
		}
	}
	var s string
	if s, _, status = looksLikeNumber(p, true); status.IsOK() {
		defer func() { p.Stop(status.Start) }()
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		result, status.Err = bwstr.ParseUint(s)
		status.UnexpectedIfErr(p)
	}
	return
}

// ============================================================================

func ParseNumber(p bwparse.I, optOpt ...Opt) (result *bwtype.Number, status Status) {
	opt := getOpt(optOpt)
	result, status = parseNumber(p, opt, RangeLimitNone)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

// ============================================================================

// func ArrayOfString(p bwparse.I, optOpt ...Opt) (result []string, status Status) {
//  opt := getOpt(optOpt)
//  return parseArrayOfString(p, opt, false)
// }

// ============================================================================

func ParseArray(p bwparse.I, optOpt ...Opt) (result []interface{}, status Status) {
	opt := getOpt(optOpt)
	if status = parseDelimitedOptionalCommaSeparated(p, '[', ']', opt, func(on On, base bw.ValPath) (err error) {
		if result == nil {
			result = []interface{}{}
			on.Opt.Path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
		}
		if err == nil {
			var vals []interface{}
			var st Status
			if vals, st = parseArrayOfString(p, *on.Opt, true); st.Err == nil {
				if st.OK {
					result = append(result, vals...)
				} else {
					if opt.OnParseArrayElem != nil {
						var newResult []interface{}
						if newResult, st = opt.OnParseArrayElem(on, result); st.IsOK() {
							result = newResult
						}
					}
					if st.Err == nil && !st.OK {
						var val interface{}
						if val, st = ParseVal(p, opt); st.IsOK() {
							result = append(result, val)
						}
					}
				}
				on.Opt.Path[len(on.Opt.Path)-1].Idx = len(result)
			}
			err = st.Err
		}
		return
	}); status.IsOK() {
		if result == nil {
			result = []interface{}{}
		}
	}
	return
}

// ============================================================================

func ParseMap(p bwparse.I, optOpt ...Opt) (result map[string]interface{}, status Status) {
	opt := getOpt(optOpt)
	var path bw.ValPath
	if status = parseDelimitedOptionalCommaSeparated(p, '{', '}', opt, func(on On, base bw.ValPath) error {
		if result == nil {
			result = map[string]interface{}{}
			path = append(base, bw.ValPathItem{Type: bw.ValPathItemKey})
		}
		var (
			key string
		)
		onKey := func(s string, start *bwparse.Start) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				on.Opt.Path = base
				on.Start = start
				if _, b := result[key]; b {
					err = p.Error(bwparse.E{
						Start: start,
						Fmt:   bw.Fmt(ansi.String("duplicate key <ansiErr>%s<ansi>"), start.Suffix()),
					})
				} else {
					err = opt.OnValidateMapKey(on, bwmap.M(result), key)
				}
			}
			return
		}
		var st Status
		if st = processOn(p,
			onString{opt: opt, f: onKey},
			onId{opt: opt, f: onKey},
		); st.IsOK() {
			var isSpaceSkipped bool
			if isSpaceSkipped, st.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); st.Err == nil {

				var isSeparatorSkipped bool
				if bwparse.SkipRunes(p, ':') {
					isSeparatorSkipped = true
				} else if bwparse.SkipRunes(p, '=') {
					if bwparse.SkipRunes(p, '>') {
						isSeparatorSkipped = true
					} else {
						return bwparse.Expects(p, nil, fmt.Sprintf(ansiVal, string('>')))
					}
				}

				if st.Err == nil {
					if isSeparatorSkipped {
						_, st.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF)
					}
					if !(isSpaceSkipped || isSeparatorSkipped) {
						st.Err = bwparse.ExpectsSpace(p)
					} else {
						path[len(path)-1].Name = key
						on.Opt.Path = path
						on.Start = p.Start()
						defer func() { p.Stop(on.Start) }()

						st.OK = false
						if opt.OnParseMapElem != nil {
							st = opt.OnParseMapElem(on, bwmap.M(result), key)
						}
						if st.Err == nil && !st.OK {
							result[key], st = ParseVal(p, opt)
						}
					}
				}
			}
		}
		if !st.OK && st.Err == nil {
			st.Err = bwparse.Unexpected(p)
		}
		return st.Err
	}); status.IsOK() {
		if result == nil {
			result = map[string]interface{}{}
		}
	}
	return
}

// ============================================================================

func ParseOrderedMap(p bwparse.I, optOpt ...Opt) (result *bwmap.Ordered, status Status) {
	opt := getOpt(optOpt)
	var path bw.ValPath
	if status = parseDelimitedOptionalCommaSeparated(p, '{', '}', opt, func(on On, base bw.ValPath) error {
		if result == nil {
			result = bwmap.OrderedNew()
			path = append(base, bw.ValPathItem{Type: bw.ValPathItemKey})
		}
		var (
			key string
		)
		onKey := func(s string, start *bwparse.Start) (err error) {
			key = s
			if opt.OnValidateMapKey != nil {
				on.Opt.Path = base
				on.Start = start
				if _, b := result.Get(key); b {
					err = p.Error(bwparse.E{
						Start: start,
						Fmt:   bw.Fmt(ansi.String("duplicate key <ansiErr>%s<ansi>"), start.Suffix()),
					})
				} else {
					err = opt.OnValidateMapKey(on, result, key)
				}
			}
			return
		}
		var st Status
		if st = processOn(p,
			onString{opt: opt, f: onKey},
			onId{opt: opt, f: onKey},
		); st.IsOK() {
			var isSpaceSkipped bool
			if isSpaceSkipped, st.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); st.Err == nil {

				var isSeparatorSkipped bool
				if bwparse.SkipRunes(p, ':') {
					isSeparatorSkipped = true
				} else if bwparse.SkipRunes(p, '=') {
					if bwparse.SkipRunes(p, '>') {
						isSeparatorSkipped = true
					} else {
						return bwparse.Expects(p, nil, fmt.Sprintf(ansiVal, string('>')))
					}
				}

				if st.Err == nil {
					if isSeparatorSkipped {
						_, st.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF)
					}
					if !(isSpaceSkipped || isSeparatorSkipped) {
						st.Err = bwparse.ExpectsSpace(p)
					} else {
						path[len(path)-1].Name = key
						on.Opt.Path = path
						on.Start = p.Start()
						defer func() { p.Stop(on.Start) }()

						st.OK = false
						if opt.OnParseMapElem != nil {
							st = opt.OnParseMapElem(on, result, key)
						}
						if st.Err == nil && !st.OK {
							var val interface{}
							val, st = ParseVal(p, opt)
							if st.IsOK() {
								result.Set(key, val)
							}
						}
					}
				}
			}
		}
		if !st.OK && st.Err == nil {
			st.Err = bwparse.Unexpected(p)
		}
		return st.Err
	}); status.IsOK() {
		if result == nil {
			result = bwmap.OrderedNew()
		}
	}
	return
}

// ============================================================================

type PathA struct {
	Bases     []bw.ValPath
	isSubPath bool
}

type PathOpt struct {
	Opt
	Bases     []bw.ValPath
	isSubPath bool
}

func ParsePath(p bwparse.I, optOpt ...PathOpt) (result bw.ValPath, status Status) {
	opt := getPathOpt(optOpt)
	result, status = parsePath(p, opt)
	if status.OK {
		p.Stop(status.Start)
	}
	return
}

func PathFrom(s string, bases ...bw.ValPath) (result bw.ValPath, err error) {
	var p *bwparse.P
	if p, err = bwparse.From(bwrune.S{s}); err != nil {
		return
	}
	if result, err = ParsePathContent(p, PathOpt{Bases: bases}); err == nil {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	return
}

func MustPathFrom(s string, bases ...bw.ValPath) (result bw.ValPath) {
	var err error
	if result, err = PathFrom(s, bases...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func ParsePathContent(p bwparse.I, optOpt ...PathOpt) (bw.ValPath, error) {
	opt := getPathOpt(optOpt)
	p.Forward(bwparse.Initial)
	var (
		vpi           bw.ValPathItem
		isEmptyResult bool
		st            Status
	)
	result := bw.ValPath{}
	strictId := func(opt Opt) Opt {
		opt.StrictId = true
		return opt
	}
	for st.Err == nil {
		isEmptyResult = len(result) == 0
		if isEmptyResult && p.Curr().Rune() == '.' {
			p.Forward(1)
			if len(opt.Bases) > 0 {
				result = append(result, opt.Bases[0]...)
			} else {
				break
			}
		} else if st = processOn(p,
			onInt{opt: opt.Opt, f: func(idx int, start *bwparse.Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemIdx, Idx: idx}
				return
			}},
			onString{opt: opt.Opt, f: func(s string, start *bwparse.Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Name: s}
				return
			}},
			onId{opt: strictId(opt.Opt), f: func(s string, start *bwparse.Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemKey, Name: s}
				return
			}},
			onSubPath{opt: opt, f: func(path bw.ValPath, start *bwparse.Start) (err error) {
				vpi = bw.ValPathItem{Type: bw.ValPathItemPath, Path: path}
				return
			}},
		); st.Err == nil {
			if st.OK {
				result = append(result, vpi)
				// } else if bwparse.SkipRunes(p, '#') {
				//  result = append(result, bw.ValPathItem{Type: bw.ValPathItemHash})
				//  break
			} else if isEmptyResult && bwparse.SkipRunes(p, '$') {
				st = processOn(p,
					onInt{opt: opt.Opt, f: func(idx int, start *bwparse.Start) (err error) {
						l := len(opt.Bases)
						if nidx, b := bw.NormalIdx(idx, l); b {
							result = append(result, opt.Bases[nidx]...)
						} else {
							err = p.Error(bwparse.E{start, bw.Fmt(ansiUnexpectedBasePathIdx, idx, l)})
						}
						return
					}},
					onId{opt: strictId(opt.Opt), f: func(s string, start *bwparse.Start) (err error) {
						result = append(result, bw.ValPathItem{Type: bw.ValPathItemVar, Name: s})
						return
					}},
					onString{opt: opt.Opt, f: func(s string, start *bwparse.Start) (err error) {
						result = append(result, bw.ValPathItem{Type: bw.ValPathItemVar, Name: s})
						return
					}},
				)
			}
			if st.Err == nil {
				if !st.OK {
					st.Err = bwparse.Unexpected(p)
				} else {
					if !opt.isSubPath && bwparse.SkipRunes(p, '?') {
						result[len(result)-1].IsOptional = true
					}
					if bwparse.CanSkipRunes(p, '.', '.') || !bwparse.SkipRunes(p, '.') {
						break
					}
				}
			}
		}
	}
	return result, st.Err
}

// ============================================================================

func ParseRange(p bwparse.I, optOpt ...Opt) (result *bwtype.Range, status Status) {
	opt := getOpt(optOpt)

	var (
		min, rangeLimitVal interface{}
		isNumber           bool
		isPath             bool
		justParsedPath     bw.ValPath
	)

	hasKind := func(kind bwtype.ValKind, rlk RangeLimitKind) (result bool) {
		var ks bwtype.ValKindSet
		if rlk == RangeLimitMin {
			ks = opt.RangeLimitMinOwnKindSet
		} else {
			ks = opt.RangeLimitMaxOwnKindSet
		}
		if len(ks) != 0 {
			result = ks.Has(kind)
		} else if len(opt.KindSet) == 0 {
			result = true
		} else if !opt.ExcludeKinds {
			result = opt.KindSet.Has(kind)
		} else if opt.ExcludeKinds {
			result = !opt.KindSet.Has(kind)
		}
		return
	}

	onArgs := func(rlk RangeLimitKind) (onArgs []on) {
		if hasKind(bwtype.ValPath, rlk) {
			onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *bwparse.Start) (err error) {
				justParsedPath = path
				rangeLimitVal = path
				isPath = true
				return
			}})
		}

		if hasKind(bwtype.ValNumber, rlk) {
			onArgs = append(onArgs, onNumber{opt: opt, rlk: rlk, f: func(n *bwtype.Number, start *bwparse.Start) (err error) {
				rangeLimitVal = n
				isNumber = true
				return
			}})
		} else if hasKind(bwtype.ValInt, rlk) {
			onArgs = append(onArgs, onInt{opt: opt, rlk: rlk, f: func(i int, start *bwparse.Start) (err error) {
				rangeLimitVal = i
				isNumber = true
				return
			}})
		} else if hasKind(bwtype.ValUint, rlk) {
			onArgs = append(onArgs, onUint{opt: opt, f: func(u uint, start *bwparse.Start) (err error) {
				rangeLimitVal = u
				isNumber = true
				return
			}})
		}
		return
	}

	pp := bwparse.ProxyFrom(p)
	{
		if status = processOn(pp, onArgs(RangeLimitMin)...); status.Err != nil {
			if status.OK {
				pp.Stop(status.Start)
			}
			return
		}
		min = rangeLimitVal

		if status.OK = bwparse.CanSkipRunes(pp, '.', '.'); !status.OK {
			if isNumber || isPath {
				pp.Stop(status.Start)
				ps := status.Start.PosInfo()
				if isNumber {
					ps.JustParsed = numberResult{bwtype.MustNumberFrom(min), status.Start}
				} else if isPath {
					ps.JustParsed = pathResult{justParsedPath, status.Start}
				}
			}
			status = Status{}
			return
		}
	}
	status.Start = p.Start()
	defer func() { p.Stop(status.Start) }()
	p.Forward(pp.Ofs() + 2)
	pp = nil

	rangeLimitVal = nil
	var st Status
	if st = processOn(p, onArgs(RangeLimitMax)...); st.OK {
		p.Stop(st.Start)
	}
	if st.Err != nil {
		status.Err = st.Err
	} else {
		result, status.Err = bwtype.RangeFrom(bwtype.A{Min: min, Max: rangeLimitVal})
	}

	return
}

// ============================================================================

func ParseVal(p bwparse.I, optOpt ...Opt) (result interface{}, status Status) {
	opt := getOpt(optOpt)
	var onArgs []on
	kindSet := bwtype.ValKindSet{}
	kinds := []bwtype.ValKind{}
	kindSetIsEmpty := len(opt.KindSet) == 0
	hasKind := func(kind bwtype.ValKind) (result bool) {
		if kindSetIsEmpty {
			result = true
		} else if !opt.ExcludeKinds {
			result = opt.KindSet.Has(kind)
		} else if opt.ExcludeKinds {
			result = !opt.KindSet.Has(kind)
		}
		if result {
			if !kindSet.Has(kind) {
				kinds = append(kinds, kind)
				kindSet.Add(kind)
			}
		}
		return
	}
	if hasKind(bwtype.ValArray) {
		onArgs = append(onArgs, onArray{opt: opt, f: func(vals []interface{}, start *bwparse.Start) (err error) {
			if opt.OnValidateArray != nil {
				err = opt.OnValidateArray(On{p, start, &opt}, vals)
			}
			if err == nil {
				result = vals
			}
			return
		}})
		onArgs = append(onArgs, onArrayOfString{opt: opt, f: func(vals []interface{}, start *bwparse.Start) (err error) {
			if opt.OnValidateArray != nil {
				err = opt.OnValidateArray(On{p, start, &opt}, vals)
			}
			if err == nil {
				result = vals
			}
			return
		}})
	}
	if hasKind(bwtype.ValString) {
		onArgs = append(onArgs, onString{opt: opt, f: func(s string, start *bwparse.Start) (err error) {
			if opt.OnValidateString != nil {
				err = opt.OnValidateString(On{p, start, &opt}, s)
			}
			if err == nil {
				result = s
			}
			return
		}})
	}
	if hasKind(bwtype.ValRange) {
		onArgs = append(onArgs, onRange{opt: opt, f: func(rng *bwtype.Range, start *bwparse.Start) (err error) {
			if opt.OnValidateRange != nil {
				err = opt.OnValidateRange(On{p, start, &opt}, rng)
			}
			if err == nil {
				result = rng
			}
			return
		}})
	}
	if hasKind(bwtype.ValPath) {
		onArgs = append(onArgs, onPath{opt: PathOpt{Opt: opt}, f: func(path bw.ValPath, start *bwparse.Start) (err error) {
			if opt.OnValidatePath != nil {
				err = opt.OnValidatePath(On{p, start, &opt}, path)
			}
			if err == nil {
				result = path
			}
			return
		}})
	}

	if hasKind(bwtype.ValOrderedMap) {
		onArgs = append(onArgs, onOrderedMap{opt: opt, f: func(m *bwmap.Ordered, start *bwparse.Start) (err error) {
			if opt.OnValidateMap != nil {
				err = opt.OnValidateMap(On{p, start, &opt}, m)
			}
			if err == nil {
				result = m
			}
			return
		}})
	} else if hasKind(bwtype.ValMap) {
		onArgs = append(onArgs, onMap{opt: opt, f: func(m map[string]interface{}, start *bwparse.Start) (err error) {
			if opt.OnValidateMap != nil {
				err = opt.OnValidateMap(On{p, start, &opt}, bwmap.M(m))
			}
			if err == nil {
				result = m
			}
			return
		}})
	}

	if hasKind(bwtype.ValNumber) {
		onArgs = append(onArgs, onNumber{opt: opt, f: func(n *bwtype.Number, start *bwparse.Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, n)
			}
			if err == nil {
				val := n.Val()
				if i, b := bwtype.Int(val); b {
					result = i
				} else if u, b := bwtype.Uint(val); b {
					result = u
				} else {
					result = val
				}
			}
			return
		}})
	} else if hasKind(bwtype.ValInt) {
		onArgs = append(onArgs, onInt{opt: opt, f: func(i int, start *bwparse.Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(i))
			}
			if err == nil {
				result = i
			}
			return
		}})
	} else if hasKind(bwtype.ValUint) {
		onArgs = append(onArgs, onUint{opt: opt, f: func(u uint, start *bwparse.Start) (err error) {
			if opt.OnValidateNumber != nil {
				err = opt.OnValidateNumber(On{p, start, &opt}, bwtype.MustNumberFrom(u))
			}
			if err == nil {
				result = u
			}
			return
		}})
	}
	if hasKind(bwtype.ValNil) {
		onArgs = append(onArgs, onNil{opt: opt, f: func(start *bwparse.Start) (err error) { return }})
	}
	if hasKind(bwtype.ValBool) {
		onArgs = append(onArgs, onBool{opt: opt, f: func(b bool, start *bwparse.Start) (err error) { result = b; return }})
	}
	if hasKind(bwtype.ValDef) {
		onArgs = append(onArgs, onDef{opt: opt, f: func(b bool, start *bwparse.Start) (err error) { result = b; return }})
	}
	if (len(opt.IdVals) > 0 || opt.OnId != nil) && hasKind(bwtype.ValId) || !opt.StrictId && hasKind(bwtype.ValString) {
		onArgs = append(onArgs,
			onId{opt: opt, f: func(s string, start *bwparse.Start) (err error) {
				var b bool
				if result, b = opt.IdVals[s]; !b {
					if opt.OnId != nil {
						result, b, err = opt.OnId(On{p, start, &opt}, s)
					}
				}
				if !b && err == nil {
					if opt.StrictId {
						err = p.Error(bwparse.E{start, bw.Fmt(ansiUnexpectedWord, s)})
						if expects := getIdExpects(opt, ""); len(expects) > 0 {
							err = bwparse.Expects(p, err, expects)
						}
					} else {
						if opt.OnValidateString != nil {
							err = opt.OnValidateString(On{p, start, &opt}, s)
						}
						if err == nil {
							result = s
						}
					}
				}
				return
			}},
		)
	}
	if status = processOn(p, onArgs...); !status.OK && !opt.ValFalse {
		var expects []string
		asType := func(kind bwtype.ValKind) (result string) {
			s := kind.String()
			switch kind {
			case bwtype.ValNumber, bwtype.ValInt:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitNone) {
					s = "NonNegative" + s
				}
			case bwtype.ValRange:
				if opt.NonNegativeNumber != nil && opt.NonNegativeNumber(RangeLimitMin) {
					s = s + "(Min: NonNegative)"
				}
			}
			result = fmt.Sprintf(ansiType, s)
			if kind == bwtype.ValId {
				if expects := getIdExpects(opt, "  "); len(expects) > 0 {
					result += "(" + expects + ")"
				}
			}
			return
		}
		addExpects := func(kind bwtype.ValKind) {
			expects = append(expects, asType(kind))
		}
		for _, kind := range kinds {
			addExpects(kind)
		}
		status.Err = bwparse.Expects(p, status.Err,
			bwstr.SmartJoin(bwstr.A{
				Source: bwstr.SS{
					SS: expects,
				},
				MaxLen:              60,
				NoJoinerOnMutliline: true,
			}),
		)
	}
	return
}

// ============================================================================

var (
	// ansiOptKeyIsNotOfType    string
	// ansiOptHasUnexpectedKeys string

	// ansiOK  string
	// ansiErr string

	// ansiPos                   string
	// ansiLineCol               string
	// // ansiGetSuffixAssert       string
	// // ansiUnexpectedEOF         string
	// // ansiUnexpectedChar        string
	ansiUnexpectedWord string
	// ansiOutOfRange            string
	ansiUnexpectedBasePathIdx string
	// // ansiType                  string
	ansiVal string

// ansiVarSpace              string
)

func init() {
	// ansiOptKeyIsNotOfType = ansi.String("<ansiVar>opt.%s<ansi> (<ansiVal>%#v<ansi>) is not <ansiType>%s")
	// ansiOptHasUnexpectedKeys = ansi.String("<ansiVar>opt<ansi> (<ansiVal>%s<ansi>) has unexpected keys <ansiVal>%s")

	// ansiOK = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorGreen, Bright: false})).String()
	// ansiErr = ansi.CSIFromSGRCodes(ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorRed, Bright: true})).String()

	// ansiPos = ansi.String(" at pos <ansiPath>%d<ansi>")
	// ansiLineCol = ansi.String(" at line <ansiPath>%d<ansi>, col <ansiPath>%d<ansi> (pos <ansiPath>%d<ansi>)")
	// ansiGetSuffixAssert = ansi.String("<ansiVar>ps.pos<ansi> (<ansiVal>%d<ansi>) > <ansiVar>p.curr.pos<ansi> (<ansiVal>%d<ansi>)")
	// ansiUnexpectedEOF = ansi.String("unexpected end of string")
	// ansiUnexpectedChar = ansi.String("unexpected char <ansiVal>%q<ansiReset> (<ansiVar>charCode<ansi>: <ansiVal>%d<ansi>)")
	ansiUnexpectedWord = ansi.String("unexpected `<ansiErr>%s<ansi>`")
	ansiUnexpectedBasePathIdx = ansi.String("unexpected base path idx <ansiVal>%d<ansi> (len(bases): <ansiVal>%d)")
	// ansiType = ansi.String("<ansiType>%s")
	ansiVal = ansi.String("<ansiVal>%s")
	// ansiVarSpace = ansi.String("<ansiVar>Space")
}

// func optKeyUint(opt map[string]interface{}, key string, keys *bwset.String) (result uint, ok bool) {
//   var val interface{}
//   keys.Add(key)
//   if val, ok = opt[key]; ok && val != nil {
//     if result, ok = bwtype.Uint(val); !ok {
//       bwerr.Panic(ansiOptKeyIsNotOfType, key, val, "Uint")
//     }
//   }
//   return
// }

// ============================================================================

func getOpt(optOpt []Opt) (result Opt) {
	if len(optOpt) > 0 {
		result = optOpt[0]
	}
	return
}

func getPathOpt(optOpt []PathOpt) (result PathOpt) {
	if len(optOpt) > 0 {
		result = optOpt[0]
	}
	return
}

func isOneOfId(p bwparse.I, ss []string) (needForward uint, ok bool) {
	p.Forward(bwparse.Initial)
	for _, s := range ss {
		if ok = bwparse.CanSkipRunes(p, []rune(s)...); ok {
			u := uint(len(s))
			r := p.LookAhead(u).Rune()
			if ok = !(bw.IsLetter(r) || bw.IsDigit(r)); ok {
				needForward = u
				return
			}
		}
	}
	return
}

// ============================================================================

// func (p *P) pullRune(ps PosInfo) PosInfo {
//   runePtr := bwrune.MustPull(p.prov)
//   if !ps.isEOF {
//     if ps.pos >= 0 {
//       ps.prefix += string(ps.rune)
//     }
//     if runePtr != nil {
//       if ps.rune != '\n' {
//         ps.col++
//       } else {
//         ps.line++
//         ps.col = 1
//         if int(ps.line) > int(p.preLineCount)+1 {
//           i := strings.Index(ps.prefix, "\n")
//           ps.prefix = ps.prefix[i+1:]
//           ps.prefixStart += i + 1
//         }
//       }
//     }
//     if runePtr == nil {
//       ps.rune, ps.isEOF = '\000', true
//     } else {
//       ps.rune, ps.isEOF = *runePtr, false
//     }
//     ps.pos++
//   }
//   return ps
// }

// func suffix(p bwparse.I, start Start, postLineCount uint) (suffix string) {
//   if start.ps.pos > p.Curr().pos {
//     bwerr.Panic(ansiGetSuffixAssert, start.ps.pos, p.Curr().pos)
//   }

//   var separator string
//   if p.Curr().line > 1 {
//     suffix += fmt.Sprintf(ansiLineCol, start.ps.line, start.ps.col, start.ps.pos)
//     separator = "\n"
//   } else {
//     suffix += fmt.Sprintf(ansiPos, start.ps.pos)
//     separator = " "
//   }
//   if fileSpec := p.FileSpec(); fileSpec != "" {
//     suffix += fmt.Sprintf(" of <ansiPath>%s<ansi>", bwos.ShortenFileSpec(fileSpec))
//   }
//   suffix += ":" + separator + ansiOK + start.ps.prefix

//   var needPostLines, noNeedNewline bool
//   if start.ps.pos < p.Curr().pos {
//     suffix += ansiErr + start.suffix + ansi.Reset()
//     needPostLines = true
//   } else if !p.Curr().isEOF {
//     suffix += ansiErr + string(p.Curr().rune) + ansi.Reset()
//     p.Forward(1)
//     needPostLines = true
//   }
//   noNeedNewline = p.Curr().rune == '\n'

//   i := int(postLineCount)
//   for needPostLines && !p.Curr().isEOF && i >= 0 {
//     suffix += string(p.Curr().rune)
//     if noNeedNewline = p.Curr().rune == '\n'; noNeedNewline {
//       i -= 1
//     }
//     p.Forward(1)
//   }

//   if !noNeedNewline {
//     suffix += string('\n')
//   }
//   return
// }

// func (p *P) forward() {
//   if !p.curr.isEOF {
//     for _, start := range p.starts {
//       start.suffix += string(p.curr.rune)
//     }
//   }
//   if len(p.next) == 0 {
//     newCurr := p.pullRune(*p.curr)
//     p.curr = &newCurr
//   } else {
//     last := len(p.next) - 1
//     p.curr, p.next = p.next[last], p.next[:last]
//   }
// }

// // ============================================================================

// type bwparse.Proxy struct {
//   p      bwparse.I
//   ofs    uint
//   starts map[int]*bwparse.Start
// }

// func (p *bwparse.Proxy) FileSpec() string {
//   return p.p.FileSpec()
// }

// func (p *bwparse.Proxy) Close() error {
//   return p.p.Close()
// }

// func (p *bwparse.Proxy) Curr() *PosInfo {
//   result := p.p.LookAhead(p.ofs)
//   return result
// }

// func (p *bwparse.Proxy) Forward(count uint) {
//   if count == 0 {
//     p.p.Forward(0)
//   } else {
//     ps := p.Curr()
//     for !ps.isEOF && count > 0 {
//       for _, start := range p.starts {
//         start.suffix += string(ps.rune)
//       }
//       count--
//       p.ofs++
//       ps = p.Curr()
//     }
//   }
// }

// func (p *bwparse.Proxy) LookAhead(ofs uint) *PosInfo {
//   return p.p.LookAhead(p.ofs + ofs)
// }

// func (p *bwparse.Proxy) Error(a E) error {
//   p.p.Forward(p.ofs)
//   return p.p.Error(a)
// }

// func (p *bwparse.Proxy) Start() (result *bwparse.Start) {
//   p.p.Forward(bwparse.Initial)
//   curr := p.Curr()
//   var ok bool
//   if result, ok = p.starts[curr.pos]; !ok {
//     result = &Start{ps: curr}
//     if p.starts == nil {
//       p.starts = map[int]*bwparse.Start{}
//     }
//     p.starts[curr.pos] = result
//   }
//   return
// }

// ============================================================================

type on interface {
	IsOn()
}

type onInt struct {
	f   func(i int, start *bwparse.Start) (err error)
	opt Opt
	rlk RangeLimitKind
}

func (onInt) IsOn() {}

type onUint struct {
	f   func(u uint, start *bwparse.Start) (err error)
	opt Opt
}

func (onUint) IsOn() {}

type onNumber struct {
	f   func(n *bwtype.Number, start *bwparse.Start) (err error)
	opt Opt
	rlk RangeLimitKind
}

func (onNumber) IsOn() {}

type onRange struct {
	f   func(rng *bwtype.Range, start *bwparse.Start) (err error)
	opt Opt
}

func (onRange) IsOn() {}

type onId struct {
	f   func(s string, start *bwparse.Start) (err error)
	opt Opt
}

func (onId) IsOn() {}

type onString struct {
	f   func(s string, start *bwparse.Start) (err error)
	opt Opt
}

func (onString) IsOn() {}

type onSubPath struct {
	f   func(path bw.ValPath, start *bwparse.Start) (err error)
	opt PathOpt
}

func (onSubPath) IsOn() {}

type onPath struct {
	f   func(path bw.ValPath, start *bwparse.Start) (err error)
	opt PathOpt
}

func (onPath) IsOn() {}

type onArray struct {
	f   func(vals []interface{}, start *bwparse.Start) (err error)
	opt Opt
}

func (onArray) IsOn() {}

type onArrayOfString struct {
	f   func(vals []interface{}, start *bwparse.Start) (err error)
	opt Opt
}

func (onArrayOfString) IsOn() {}

type onMap struct {
	f   func(m map[string]interface{}, start *bwparse.Start) (err error)
	opt Opt
}

func (onMap) IsOn() {}

type onOrderedMap struct {
	f   func(m *bwmap.Ordered, start *bwparse.Start) (err error)
	opt Opt
}

func (onOrderedMap) IsOn() {}

type onNil struct {
	f   func(start *bwparse.Start) (err error)
	opt Opt
}

func (onNil) IsOn() {}

type onBool struct {
	f   func(b bool, start *bwparse.Start) (err error)
	opt Opt
}

func (onBool) IsOn() {}

type onDef struct {
	f   func(def *bwtype.Def, start *bwparse.Start) (err error)
	opt bw.ValPathProvider
}

func (onDef) IsOn() {}

// ============================================================================

func processOn(p bwparse.I, processors ...on) (status Status) {
	var (
		i    int
		u    uint
		n    *bwtype.Number
		s    string
		path bw.ValPath
		vals []interface{}
		m    map[string]interface{}
		o    *bwmap.Ordered
		b    bool
		rng  *bwtype.Range
		def  *bwtype.Def
	)
	p.Forward(bwparse.Initial)
	for _, processor := range processors {
		switch t := processor.(type) {
		case onInt:
			i, status = parseInt(p, t.opt, t.rlk)
		case onUint:
			u, status = ParseUint(p, t.opt)
		case onNumber:
			n, status = parseNumber(p, t.opt, t.rlk)
		case onRange:
			rng, status = ParseRange(p, t.opt)
		case onString:
			s, status = ParseString(p, t.opt)
		case onId:
			s, status = ParseId(p, t.opt)
		case onSubPath:
			path, status = parseSubPath(p, t.opt)
		case onPath:
			path, status = ParsePath(p, t.opt)
		case onArray:
			vals, status = ParseArray(p, t.opt)
		case onArrayOfString:
			vals, status = parseArrayOfString(p, t.opt, false)
		case onMap:
			m, status = ParseMap(p, t.opt)
		case onOrderedMap:
			o, status = ParseOrderedMap(p, t.opt)
		case onNil:
			status = ParseNil(p, t.opt)
		case onBool:
			b, status = ParseBool(p, t.opt)
		case onDef:
			def, status = ParseDef(p, t.opt)
		}
		if status.Err != nil {
			return
		}
		if status.OK {
			switch t := processor.(type) {
			case onInt:
				status.Err = t.f(i, status.Start)
			case onUint:
				status.Err = t.f(u, status.Start)
			case onNumber:
				status.Err = t.f(n, status.Start)
			case onRange:
				status.Err = t.f(rng, status.Start)
			case onString:
				status.Err = t.f(s, status.Start)
			case onId:
				status.Err = t.f(s, status.Start)
			case onSubPath:
				status.Err = t.f(path, status.Start)
			case onPath:
				status.Err = t.f(path, status.Start)
			case onArray:
				status.Err = t.f(vals, status.Start)
			case onArrayOfString:
				status.Err = t.f(vals, status.Start)
			case onMap:
				status.Err = t.f(m, status.Start)
			case onOrderedMap:
				status.Err = t.f(o, status.Start)
			case onNil:
				status.Err = t.f(status.Start)
			case onBool:
				status.Err = t.f(b, status.Start)
			case onDef:
				status.Err = t.f(def, status.Start)
			}
			return
		}
	}
	return
}

// ============================================================================

func parseArrayOfString(p bwparse.I, opt Opt, isEmbeded bool) (result []interface{}, status Status) {
	var needForward uint
	if status.OK = p.Curr().Rune() == '<'; !status.OK {
		if status.OK = bwparse.CanSkipRunes(p, 'q', 'w') && bw.IsPunctOrSymbol(p.LookAhead(2).Rune()); !status.OK {
			return
		}
		needForward = 2
	}
	status.Start = p.Start()
	defer func() { p.Stop(status.Start) }()
	p.Forward(needForward)

	delimiter := p.Curr().Rune()
	if r, b := bw.Braces()[delimiter]; b {
		delimiter = r
	}
	p.Forward(1)
	ss := []string{}
	on := On{p, status.Start, &opt}
	if !isEmbeded {
		base := opt.Path
		on.Opt.Path = append(base, bw.ValPathItem{Type: bw.ValPathItemIdx})
		defer func() { on.Opt.Path = base }()
	}
	parseItem := func(r rune) {
		on.Start = p.Start()
		defer func() { p.Stop(on.Start) }()
		var s string
		for status.Err == nil && !(unicode.IsSpace(r) || r == delimiter) {
			s += string(r)
			p.Forward(1)
			if status.Err = bwparse.CheckNotEOF(p); status.Err == nil {
				r = p.Curr().Rune()
			}
		}
		if status.Err == nil {
			if !isEmbeded && opt.OnValidateArrayOfStringElem != nil {
				status.Err = opt.OnValidateArrayOfStringElem(on, ss, s)
			} else if opt.OnValidateString != nil {
				status.Err = opt.OnValidateString(on, s)
			}
			if status.Err == nil {
				ss = append(ss, s)
			}
		}
		on.Opt.Path[len(on.Opt.Path)-1].Idx++
	}
	for status.Err == nil {
		if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
			r := p.Curr().Rune()
			if r == delimiter {
				p.Forward(1)
				break
			}
			parseItem(r)
		}
	}
	if status.Err == nil {
		result = []interface{}{}
		for _, s := range ss {
			result = append(result, s)
		}
	}
	return
}

// ============================================================================

func parsePath(p bwparse.I, opt PathOpt) (result bw.ValPath, status Status) {
	var justParsed pathResult
	curr := p.Curr()
	if justParsed, status.OK = curr.JustParsed.(pathResult); status.OK {
		result = justParsed.path
		status.Start = justParsed.start
		p.Forward(uint(len(justParsed.start.Suffix())))
		return
		// } else if status.OK = curr.rune == '$'; status.OK {
		//  status.Start = p.Start()
		//  result, status.Err = ParsePathContent(p, opt)
	} else if status.OK = bwparse.CanSkipRunes(p, '('); status.OK {
		status.Start = p.Start()
		p.Forward(2)
		if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
			if result, status.Err = ParsePathContent(p, opt); status.Err == nil {
				if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
					if !bwparse.SkipRunes(p, ')') {
						status.Err = bwparse.Unexpected(p)
					}
				}
			}
		}
	}
	return
}

// ============================================================================

func looksLikeNumber(p bwparse.I, nonNegative bool) (s string, isNegative bool, status Status) {
	var (
		r         rune
		needDigit bool
	)
	p.Forward(bwparse.Initial)
	r = p.Curr().Rune()
	if status.OK = r == '+'; status.OK {
		needDigit = true
	} else if status.OK = !nonNegative && r == '-'; status.OK {
		s = string(r)
		needDigit = true
		isNegative = true
	} else if status.OK = bw.IsDigit(r); status.OK {
		s = string(r)
	} else {
		return
	}
	status.Start = p.Start()
	p.Forward(1)
	if needDigit {
		if r = p.Curr().Rune(); !bw.IsDigit(r) {
			status.Err = bwparse.Unexpected(p)
		} else {
			p.Forward(1)
			s += string(r)
		}
	}
	return
}

// ============================================================================

type numberResult struct {
	n     *bwtype.Number
	start *bwparse.Start
}

type pathResult struct {
	path  bw.ValPath
	start *bwparse.Start
}

// ============================================================================

func parseInt(p bwparse.I, opt Opt, rangeLimitKind RangeLimitKind) (result int, status Status) {
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber(rangeLimitKind)
	}
	var justParsed numberResult
	curr := p.Curr()
	if justParsed, status.OK = curr.JustParsed.(numberResult); status.OK {
		if result, status.OK = bwtype.Int(justParsed.n.Val()); status.OK {
			if status.OK = !nonNegativeNumber || result >= 0; status.OK {
				status.Start = justParsed.start
				p.Forward(uint(len(justParsed.start.Suffix())))
			}
			return
		}
	}
	var s string
	if s, _, status = looksLikeNumber(p, nonNegativeNumber); status.IsOK() {
		b := true
		for b {
			s, b = addDigit(p, s)
		}
		result, status.Err = bwstr.ParseInt(s)
		status.UnexpectedIfErr(p)
	}
	return
}

// ============================================================================

func parseNumber(p bwparse.I, opt Opt, rangeLimitKind RangeLimitKind) (result *bwtype.Number, status Status) {
	var (
		s          string
		hasDot     bool
		b          bool
		isNegative bool
		justParsed numberResult
	)
	nonNegativeNumber := false
	if opt.NonNegativeNumber != nil {
		nonNegativeNumber = opt.NonNegativeNumber(rangeLimitKind)
	}
	curr := p.Curr()
	if justParsed, status.OK = curr.JustParsed.(numberResult); status.OK {
		if nonNegativeNumber {
			if _, status.OK = bwtype.Uint(justParsed.n.Val()); !status.OK {
				status.Err = bwparse.Unexpected(p)
				return
			}
		}
		status.Start = justParsed.start
		result = justParsed.n
		p.Forward(uint(len(justParsed.start.Suffix())))
	} else if s, isNegative, status = looksLikeNumber(p, nonNegativeNumber); status.IsOK() {
		for {
			if s, b = addDigit(p, s); !b {
				if !hasDot && bwparse.CanSkipRunes(p, dotRune) {
					pi := p.LookAhead(1)
					if bw.IsDigit(pi.Rune()) {
						p.Forward(1)
						s += string(dotRune)
						hasDot = true
					} else {
						break
					}
				} else {
					break
				}
			}
		}

		if hasDot && !zeroAfterDotRegexp.MatchString(s) {
			var f float64
			if f, status.Err = strconv.ParseFloat(s, 64); status.Err == nil {
				result = bwtype.MustNumberFrom(f)
			}
		} else {
			if pos := strings.LastIndex(s, string(dotRune)); pos >= 0 {
				s = s[:pos]
			}
			if isNegative {
				var i int
				if i, status.Err = bwstr.ParseInt(s); status.Err == nil {
					result = bwtype.MustNumberFrom(i)
				}
			} else {
				var u uint
				if u, status.Err = bwstr.ParseUint(s); status.Err == nil {
					result = bwtype.MustNumberFrom(u)
				}
			}
		}
		status.UnexpectedIfErr(p)
	}
	return
}

const dotRune = '.'

var zeroAfterDotRegexp = regexp.MustCompile(`\.0+$`)

// ============================================================================

func getIdExpects(opt Opt, indent string) (expects string) {
	if len(opt.IdVals) > 0 {
		sset := bwset.String{}
		for s := range opt.IdVals {
			sset.Add(s)
		}
		var suffix string
		if opt.OnId != nil {
			suffix = ansi.String(" or <ansiVar>custom<ansi>")
		}
		expects = bwstr.SmartJoin(bwstr.A{
			Source: bwstr.SS{
				SS: sset.ToSliceOfStrings(),
				Preformat: func(s string) string {
					return fmt.Sprintf(ansi.String("<ansiVal>%s"), s)
				},
			},
			MaxLen:              uint(60 - len(suffix)),
			NoJoinerOnMutliline: true,
			InitialIndent:       indent,
		}) + suffix
	}
	return
}

// ============================================================================

func parseDelimitedOptionalCommaSeparated(p bwparse.I, openDelimiter, closeDelimiter rune, opt Opt, fn func(on On, base bw.ValPath) error) (status Status) {
	p.Forward(bwparse.Initial)
	if status.OK = bwparse.CanSkipRunes(p, openDelimiter); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		base := opt.Path
		on := On{p, status.Start, &opt}
		defer func() { on.Opt.Path = base }()
	LOOP:
		for status.Err == nil {
			if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
			NEXT:
				if bwparse.SkipRunes(p, closeDelimiter) {
					break LOOP
				}
				if status.Err = fn(on, base); status.Err == nil {
					var isSpaceSkipped bool
					if isSpaceSkipped, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
						if bwparse.SkipRunes(p, closeDelimiter) {
							break LOOP
						}
						if !bwparse.SkipRunes(p, ',') {
							if isSpaceSkipped {
								goto NEXT
							} else {
								status.Err = bwparse.ExpectsSpace(p)
							}
						}
					}
				}
			}
		}
	}
	return
}

// ============================================================================

func addDigit(p bwparse.I, s string) (string, bool) {
	r := p.Curr().Rune()
	if bw.IsDigit(r) {
		s += string(r)
	} else if r != '_' {
		return s, false
	}
	p.Forward(1)
	return s, true
}

// ============================================================================

func parseSubPath(p bwparse.I, opt PathOpt) (result bw.ValPath, status Status) {
	if status.OK = bwparse.CanSkipRunes(p, '('); status.OK {
		status.Start = p.Start()
		defer func() { p.Stop(status.Start) }()
		p.Forward(1)
		if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
			subOpt := opt
			subOpt.isSubPath = true
			subOpt.Bases = opt.Bases
			if result, status.Err = ParsePathContent(p, subOpt); status.Err == nil {
				if _, status.Err = bwparse.SkipSpace(p, bwparse.TillNonEOF); status.Err == nil {
					if p.Curr().Rune() == ')' {
						p.Forward(1)
					} else {
						status.Err = bwparse.Unexpected(p)
					}
				}
			}
		}
	}
	return
}

// ============================================================================
