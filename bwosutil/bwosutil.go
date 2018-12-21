package bwosutil

import (
	"bufio"
	"os"

	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwparse"
	"github.com/baza-winner/bwcore/bwrune"
)

func CreateFileFromTemplate(targetFileSpec, templateFileSpec string, vars map[string]string) (err error) {
	var file *os.File
	if file, err = os.Create(targetFileSpec); err != nil {
		return
	}
	defer func() {
		e := file.Close()
		if err == nil {
			err = e
		}
	}()
	w := bufio.NewWriter(file)
	var p *bwparse.P
	if p, err = bwparse.From(bwrune.F{S: templateFileSpec}); err != nil {
		return
	}
	p.Forward(bwparse.Initial)

	for !p.Curr().IsEOF() {
		for !p.Curr().IsEOF() && (p.Curr().Rune() != '$' || p.LookAhead(1).Rune() != '{') {
			w.WriteRune(p.Curr().Rune())
			p.Forward(1)
		}
		if !p.Curr().IsEOF() {
			p.Forward(2)
			start := p.Start()
			defer p.Stop(start)
			var id string
			for !p.Curr().IsEOF() && p.Curr().Rune() != '}' {
				id += string(p.Curr().Rune())
				p.Forward(1)
			}
			if p.Curr().IsEOF() {
				err = bwparse.Unexpected(p)
				return
			}
			if s, ok := vars[id]; ok {
				w.WriteString(s)
			} else {
				err = p.Error(bwparse.E{Start: start, Fmt: bw.Fmt("unexpected var <ansiVar>%s<ansi>", id)})
				return
			}
			p.Forward(1)
		}
	}
	if err = w.Flush(); err != nil {
		return
	}
	return
}
