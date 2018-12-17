// Package bwval реализует интерфейc bw.Val и утилиты для работы с этим интерфейсом.
package bwval

import (
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwjson"
	"github.com/baza-winner/bwcore/bwos"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwtype"
)

// ============================================================================

type PathS struct {
	S     string
	Bases []bw.ValPath
}

func (v PathS) Path() (result bw.ValPath, err error) {
	var p *bwparse.P
	if p, err = bwparse.From(bwrune.S{v.S}); err != nil {
		return
	}
	defer p.Close()
	opt := bwparse.PathOpt{Bases: v.Bases}
	if result, err = bwparse.PathContent(p, opt); err == nil {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	return
}

// ============================================================================

type PathSS struct {
	SS    []string
	Bases []bw.ValPath
}

func (v PathSS) Path() (result bw.ValPath, err error) {
	for _, s := range v.SS {
		result = append(result, bw.ValPathItem{Type: bw.ValPathItemKey, Key: s})
	}
	return
}

// ============================================================================

func MustPath(pathProvider bw.ValPathProvider) (result bw.ValPath) {
	var err error
	if result, err = pathProvider.Path(); err != nil {
		bwerr.PanicErr(bwerr.Refine(err, "invalid path: {Error}"))
	}
	return
}

// ============================================================================

func DefFrom(pp bwrune.ProviderProvider, optBaseProvider ...bw.ValPathProvider) (result *Def, err error) {
	var p *bwparse.P
	if p, err = bwparse.From(pp); err != nil {
		return
	}
	defer p.Close()
	if _, err = bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
		err = bwerr.Refine(err, "<ansiFunc>DefFrom<ansi> failed: {Error}")
		return
	}
	var st bwparse.Status
	if result, st = ParseDef(p); !st.IsOK() {
		err = st.Err
	} else {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	return
}

func MustDefFrom(pp bwrune.ProviderProvider, optBaseProvider ...bw.ValPathProvider) (result *Def) {
	var err error
	if result, err = DefFrom(pp, optBaseProvider...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

// MustPathVal - must-обертка bw.Val.PathVal()
func MustPathVal(v bw.Val, pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	path := MustPath(pathProvider)
	if result, err = v.PathVal(path, optVars...); err != nil {
		bwerr.PanicA(bwerr.Err(bwerr.Refine(err,
			ansiMustPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		)))
	}
	return result
}

// MustSetPathVal - must-обертка bw.Val.SetPathVal()
func MustSetPathVal(val interface{}, v bw.Val, pathProvider bw.ValPathProvider, optVars ...map[string]interface{}) {
	var err error
	path := MustPath(pathProvider)
	if err = v.SetPathVal(val, path, optVars...); err != nil {
		bwerr.PanicErr(bwerr.Refine(err,
			ansiMustSetPathValFailed,
			path, bwjson.Pretty(v), varsJSON(path, optVars),
		))
	}
}

// ============================================================================

type FromProvider interface {
	getVal(def *Def) (interface{}, error)
	getPath() (bw.ValPath, bool)
}

// =====================================

type F struct {
	S    string
	Vars map[string]interface{}
}

func (v F) getVal(def *Def) (result interface{}, err error) {
	var template Template
	if template, err = TemplateFrom(bwrune.F{S: v.S}, def); err != nil {
		return
	}
	return FromTemplate(template, v.Vars)
}

func (v F) getPath() (bw.ValPath, bool) {
	return bw.ValPath{bw.ValPathItem{Type: bw.ValPathItemVar, Key: bwos.ShortenFileSpec(v.S)}}, true
}

// =====================================

type S struct {
	S string
	// Def  *Def
	Vars map[string]interface{}
}

func (v S) getVal(def *Def) (result interface{}, err error) {
	var template Template
	if template, err = TemplateFrom(bwrune.S{S: v.S}, def); err != nil {
		return
	}
	result, err = FromTemplate(template, v.Vars)
	if def != nil {
		result, err = Holder{Val: result}.ValidVal(*def)
	}
	return
}

func (v S) getPath() (result bw.ValPath, ok bool) {
	return
}

// =====================================

type T struct {
	T Template
	// Def  *Def
	Vars map[string]interface{}
}

func (v T) getVal(def *Def) (result interface{}, err error) {
	if result, err = FromTemplate(v.T, v.Vars); err != nil {
		return
	}
	if def != nil {
		result, err = Holder{Val: result}.ValidVal(*def)
	}
	return
}

func (v T) getPath() (result bw.ValPath, ok bool) {
	return
}

// =====================================

type V struct {
	Val interface{}
	// Def *Def
}

func (v V) getVal(def *Def) (result interface{}, err error) {
	if def == nil {
		result = v.Val
	} else {
		result, err = Holder{Val: v.Val}.ValidVal(*def)
	}
	return
}

func (v V) getPath() (result bw.ValPath, ok bool) {
	return
}

// =====================================

type O struct {
	PathProvider bw.ValPathProvider
	OverridePath bool
	Def          *Def
}

func From(fromProvider FromProvider, optOpt ...O) (result Holder, err error) {
	result = Holder{}
	var opt O
	if len(optOpt) > 0 {
		opt = optOpt[0]
	}
	if result.Val, err = fromProvider.getVal(opt.Def); err != nil {
		return
	}
	var ok bool
	var path bw.ValPath
	if path, ok = fromProvider.getPath(); ok && !opt.OverridePath {
		result.Pth = path
	} else if opt.PathProvider != nil {
		if path, err = opt.PathProvider.Path(); err == nil {
			result.Pth = path
		}
	}
	return
}

func MustFrom(fromProvider FromProvider, optOpt ...O) (result Holder) {
	var err error
	if result, err = From(fromProvider, optOpt...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

type Template struct {
	val interface{}
}

func TemplateFrom(pp bwrune.ProviderProvider, optDef ...*Def) (result Template, err error) {
	var val interface{}
	var p *bwparse.P
	var def *Def
	if len(optDef) > 0 {
		def = optDef[0]
	}
	if p, err = bwparse.From(pp); err != nil {
		return
	}
	defer p.Close()
	if _, err = bwparse.SkipSpace(p, bwparse.TillNonEOF); err != nil {
		err = bwerr.Refine(err, "failed to TemplateFrom: {Error}")
		return
	}
	// bwdebug.Print("def", def)
	var st bwparse.Status
	if def != nil {
		val, st = ParseValByDef(p, *def)
	} else {
		val, st = bwparse.Val(p)
	}
	if !st.IsOK() {
		err = st.Err
	} else {
		_, err = bwparse.SkipSpace(p, bwparse.TillEOF)
	}
	if err != nil {
		bwerr.PanicErr(err)
	}
	result = Template{val: val}
	return
}

func MustTemplateFrom(pp bwrune.ProviderProvider, optDef ...*Def) (result Template) {
	var err error
	if result, err = TemplateFrom(pp, optDef...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================

func FromTemplate(template Template, optVars ...map[string]interface{}) (result interface{}, err error) {
	return expandPaths(template.val, template.val, true, optVars...)
}

func MustFromTemplate(template Template, optVars ...map[string]interface{}) (result interface{}) {
	var err error
	if result, err = FromTemplate(template, optVars...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

func expandPaths(val interface{}, rootVal interface{}, isRoot bool, optVars ...map[string]interface{}) (result interface{}, err error) {
	var path bw.ValPath
	var ok bool
	if path, ok = val.(bw.ValPath); ok {
		var h Holder
		if isRoot {
			h = Holder{}
		} else {
			h = Holder{Val: rootVal}
		}
		result, err = h.PathVal(path, optVars...)
	} else {
		result = val
		switch t, kind := bwtype.Kind(val, bwtype.ValKindSetFrom(bwtype.ValMap, bwtype.ValArray)); kind {
		case bwtype.ValMap:
			m := t.(map[string]interface{})
			for key, val := range m {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				m[key] = val
			}
		case bwtype.ValArray:
			vals := t.([]interface{})
			for i, val := range vals {
				if val, err = expandPaths(val, rootVal, false, optVars...); err != nil {
					return
				}
				vals[i] = val
			}
		}
	}
	return
}

// ============================================================================
