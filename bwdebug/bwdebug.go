package bwdebug

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwerr/where"
	"github.com/baza-winner/bwcore/bwjson"
)

var (
	ansiDebugMark             string
	ansiDebugVarName          string
	ansiExpectsVarVal         string
	ansiMustBeString          string
	ansiMustBeNonEmptyString  string
	ansiMustBeNonEmptyVarName string
)

type printCase struct {
	suffixRegex *regexp.Regexp
	getString   func(arg interface{}) (result string)
	// ansiString  string
	// getArg      func(arg interface{}) interface{}
}

var printCases []printCase

func initPrintCases() {
	printCases = []printCase{
		{
			suffixRegex: regexp.MustCompile(":v$"),
			getString: func(arg interface{}) string {
				return bw.Spew.Sprintf(ansi.String("<ansiVal>%v<ansi>"), arg)
			},
		},
		{
			suffixRegex: regexp.MustCompile(":#v$"),
			getString: func(arg interface{}) string {
				return bw.Spew.Sprintf(ansi.String("<ansiVal>%#v<ansi>"), arg)
			},
		},
		{
			suffixRegex: regexp.MustCompile(":T$"),
			getString: func(arg interface{}) string {
				return fmt.Sprintf(ansi.String("<ansiVal>%T<ansi>"), arg)
			},
		},
		{
			suffixRegex: regexp.MustCompile(":s$"),
			getString: func(arg interface{}) string {
				switch t := arg.(type) {
				case rune:
					arg = string(t)
				case fmt.Stringer:
					arg = t.String()
				}
				return fmt.Sprintf(ansi.String("<ansiVal>%s<ansi>"), arg)
			},
		},
		{
			suffixRegex: regexp.MustCompile(":json$"),
			getString: func(arg interface{}) string {
				return fmt.Sprintf(ansi.String("<ansiVal>%s<ansi>"), bwjson.Pretty(arg))
			},
		},
	}
}

func init() {
	initPrintCases()
	ansi.MustAddTag("ansiDebugMark",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorYellow, Bright: true}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
	)
	ansiDebugMark = ansi.String("<ansiDebugMark>%s<ansi>, ")
	ansiDebugVarName = ansi.String("<ansiVar>%s<ansi>: ")
	ansiExpectsVarVal = ansi.String("expects val for <ansiVar>%s")
	ansiMustBeString = "<ansiVar>args<ansiPath>.%d<ansi> (<ansiVal>%#v<ansi>) must be <ansiType>string"
	ansiMustBeNonEmptyString = "<ansiVar>args<ansiPath>.%d<ansi> must be <ansiType>non empty string"
}

func Print(args ...interface{}) {
	if s, err := stringToPrint(1, args...); err != nil {
		panic(err)
	} else {
		fmt.Println(s)
	}
}

func stringToPrint(depth uint, args ...interface{}) (result string, err error) {
	markPrefix := ""
	// fmtString := ""
	// fmtArgs := []interface{}{}
	expectsVal := false
	lastVar := ""
	i := 0
	// type valueFormat uint8
	// const (
	// 	vfVerbose valueFormat = iota
	// 	vfVerboseHash
	// 	vfString
	// 	vfType
	// 	vfJSON
	// )
	var pcToUse *printCase
	pcDefault := &printCases[0]
	for _, arg := range args {
		i++
		if expectsVal == true {
			pc := pcToUse
			if pc == nil {
				pc = pcDefault
			}
			result += pc.getString(arg)
			// if pc.getArg != nil {
			// 	arg = pc.getArg(arg)
			// }
			// fmtArgs = append(fmtArgs, arg)
			expectsVal = false
		} else if valueOf := reflect.ValueOf(arg); valueOf.Kind() != reflect.String {
			err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiMustBeString, Args: bw.Args(i, arg)})
			return
		} else if s := valueOf.String(); len(s) == 0 {
			err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiMustBeNonEmptyString, Args: bw.Args(i)})
			return
		} else if s[0:1] == "!" {
			markPrefix += fmt.Sprintf(ansiDebugMark, s)
		} else {
			if len(result) > 0 {
				result += ", "
			}
			pcToUse = nil
			for i, v := range printCases {
				if s = v.suffixRegex.ReplaceAllStringFunc(s, func(string) string {
					pcToUse = &printCases[i]
					return ""
				}); pcToUse != nil {
					break
				}
			}
			if len(s) == 0 {
				pcDefault = pcToUse
			} else {
				result += fmt.Sprintf(ansiDebugVarName, s)
				lastVar = s
				expectsVal = true
			}
		}
	}
	if expectsVal {
		err = bwerr.FromA(bwerr.A{Depth: 1, Fmt: ansiExpectsVarVal, Args: bw.Args(lastVar)})
		return
	}
	result = markPrefix +
		where.MustFrom(1+depth).String() +
		": " +
		result
		// ansi.String(bw.Spew.Sprintf(fmtString, fmtArgs...))
		// ansi.String(fmt.Sprintf(fmtString, fmtArgs...))
	return
}
