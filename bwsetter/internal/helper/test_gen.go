package helper

import (
	"fmt"
	"sort"
	"strings"

	"github.com/baza-winner/bwcore/bwdebug"
	"github.com/baza-winner/bwcore/bwerr"
	. "github.com/dave/jennifer/jen"
)

// ============================================================================

const bwtestingPackage = "github.com/baza-winner/bwcore/bwtesting"

type TestCases map[string]TestCase
type TestCase struct {
	In  []interface{}
	Out []interface{}
}

type TestItem uint8

const (
	A TestItem = iota
	B
)

func (v *Helper) GenTestsFor(
	fgt FuncGenType,
	funcName string,
	pt ParamType,
	rt ReturnType,
	testData interface{},
) {
	testName := funcName
	if fgt != SimpleFunc {
		testName = v.IdSet + "." + funcName
	}
	testDict := Dict{}
	var testCase TestCase
	var testCases TestCases
	var ok bool
	testeeName := testName
	if rt == ReturnNone {
		if fgt == SimpleFunc {
			bwerr.Panic(
				"return type of <ansiPath>%s<ansi> is %s",
				testName, rt,
			)
		}
		var id string
		switch pt {
		case ParamArgs:
			id = "kk..."
		case ParamSlice:
			id = "kk"
		case ParamSet:
			id = "s"
		default:
			bwerr.Panic(
				"param type of <ansiPath>%s<ansi> is %s",
				testName, pt,
			)
		}
		testeeName = "_" + funcName + "TestHelper"
		FuncGen{v, fgt}.Func()("", testeeName, pt, ReturnSet, []*Statement{
			Id("result").Op(":=").Id("v").Dot("Copy").Call(),
			Id("result").Dot(funcName).Call(Id(id)),
			Return(Id("result")),
		}, nil)
		testeeName = v.IdSet + "." + testeeName
		rt = ReturnSet
	}
	if testCase, ok = testData.(TestCase); ok {
		testDict[Lit(testName)] = v.TestCaseValues(fgt != SimpleFunc, funcName, pt, rt, testCase)
	} else if testCases, ok = testData.(TestCases); !ok {
		bwerr.Panic(
			"<ansiPath>%s<ansi> <ansiVar>TestData<ansi> expected to be <ansiVal>%s<ansi> or <ansiVal>%s",
			testName, "TestCase", "TestCases",
		)
	} else {
		testCaseNames := []string{}
		for testCaseName, _ := range testCases {
			testCaseNames = append(testCaseNames, testCaseName)
		}
		sort.Strings(testCaseNames)
		for _, testCaseName := range testCaseNames {
			testCase := testCases[testCaseName]
			testDict[Lit(testName+": "+testCaseName)] = v.TestCaseValues(fgt != SimpleFunc, funcName, pt, rt, testCase)
		}
	}
	v.tests = append(v.tests,
		Qual(bwtestingPackage, "BwRunTests").Call(
			Id("t"),
			Id(testeeName),
			Map(String()).Qual(bwtestingPackage, "Case").Values(testDict),
		),
	)
}

func (v *Helper) TestCaseValues(
	isMethod bool,
	funcName string,
	pt ParamType,
	rt ReturnType,
	testCase TestCase,
) Code {
	testeeName := funcName
	if isMethod {
		testeeName = v.IdSet + "." + testeeName
	}
	d := Dict{}

	Params := []ParamType{}

	if isMethod {
		Params = append(Params, ParamSet)
	}
	if pt != ParamNone {
		Params = append(Params, pt)
	}
	if len(testCase.In) != len(Params) {
		bwdebug.Print("pt", pt.String())
		bwerr.Panic(
			"<ansiPath>%s<ansi> <ansiVar>testCase.In<ansi> expects to have <ansiVal>%d<ansi> item(s), but found <ansiVal>%d",
			testeeName, len(Params), len(testCase.In),
		)
	}
	pt2getValues := map[ParamType]GetValues{
		ParamSet:   v.getValuesOfSet,
		ParamArg:   v.getValuesOfArg,
		ParamArgs:  v.getValuesOfArgs,
		ParamSlice: v.getValuesOfSlice,
	}
	inValues := []Code{}
	for i, testCaseData := range testCase.In {
		var pt ParamType
		if i < len(Params) {
			pt = Params[i]
		} else {
			pt = Params[len(Params)-1]
		}
		if getValues, ok := pt2getValues[pt]; !ok {
			bwerr.Panic("pt: %s", pt)
		} else if values, err := getValues(testCaseData); err != nil {
			v.panicOnErrOfGetValues(err, testeeName, In, i)
		} else {
			inValues = append(inValues, values)
		}
	}
	d[Id("In")] = Index().Interface().Values(inValues...)

	Returns := []ReturnType{}
	if rt != ReturnNone {
		Returns = append(Returns, rt)
	}
	var expectsOutNum int
	switch rt {
	case ReturnNone:
	case ReturnJSON:
		expectsOutNum = 2
	default:
		expectsOutNum = 1
	}
	if len(testCase.Out) != expectsOutNum {
		bwerr.Panic(
			"<ansiPath>%s<ansi> <ansiVar>testCase.Out<ansi> expects to have <ansiVal>%d<ansi> item(s), but found <ansiVal>%d",
			testeeName, expectsOutNum, len(testCase.Out),
		)
	}
	rt2getValues := map[ReturnType]GetValues{
		ReturnSet:    v.getValuesOfSet,
		ReturnString: v.getValuesOfString,
		ReturnBool:   v.getValuesOfBool,
		ReturnSlice:  v.getValuesOfSlice,
	}
	outValues := []Code{}
	for i, testCaseData := range testCase.Out {
		if values, ok := testCaseData.(Code); ok {
			outValues = append(outValues, values)
		} else {
			rt := Returns[i]
			if getValues, ok := rt2getValues[rt]; !ok {
				bwerr.Panic("rt: %s", rt)
			} else if values, err := getValues(testCaseData); err == nil {
				outValues = append(outValues, values)
			} else {
				v.panicOnErrOfGetValues(err, testeeName, Out, i)
			}
		}
	}
	d[Id("Out")] = Index().Interface().Values(outValues...)
	return Values(d)
}

type GetValues func(interface{}) (Code, error)

func (v *Helper) getValuesOfSet(testCaseData interface{}) (result Code, err error) {
	if testItems, ok := testCaseData.([]TestItem); !ok {
		err = v.errOfGetValues("[]TestItem", testCaseData)
	} else {
		result = Id(v.IdSet).Values(DictFunc(func(d Dict) {
			for _, item := range testItems {
				d[Id(v.TestItemString(item))] = Struct().Values()
			}
		}))
	}
	return
}

func (v *Helper) getValuesOfSlice(testCaseData interface{}) (result Code, err error) {
	if testItems, ok := testCaseData.([]TestItem); !ok {
		err = v.errOfGetValues("[]TestItem", testCaseData)
	} else {
		values := []Code{}
		for _, item := range testItems {
			values = append(values, Id(v.TestItemString(item)))
		}
		result = Index().Id(v.IdItem).Values(values...)
	}
	return
}

func (v *Helper) getValuesOfArgs(testCaseData interface{}) (result Code, err error) {
	if testItems, ok := testCaseData.([]TestItem); !ok {
		err = v.errOfGetValues("[]TestItem", testCaseData)
	} else {
		values := []Code{}
		for _, item := range testItems {
			values = append(values, Id(v.TestItemString(item)))
		}
		result = List(values...)
		// result = Index().Id(v.IdItem).Values(values...)
	}
	return
}

func (v *Helper) getValuesOfArg(testCaseData interface{}) (result Code, err error) {
	if testItem, ok := testCaseData.(TestItem); !ok {
		err = v.errOfGetValues("TestItem", testCaseData)
	} else {
		result = Id(v.TestItemString(testItem))
	}
	return
}

func (v *Helper) getValuesOfBool(testCaseData interface{}) (result Code, err error) {
	if b, ok := testCaseData.(bool); !ok {
		err = v.errOfGetValues("bool", testCaseData)
	} else {
		if b {
			result = True()
		} else {
			result = False()
		}
	}
	return
}

func (v *Helper) getValuesOfString(testCaseData interface{}) (result Code, err error) {
	if s, ok := testCaseData.(string); !ok {
		err = v.errOfGetValues("string", testCaseData)
	} else {
		result = Lit(s)
	}
	return
}

func (v *Helper) errOfGetValues(typeName string, testCaseData interface{}) error {
	return fmt.Errorf(
		"to be <ansiVal>%s<ansi>, instead of <ansiVal>%#v",
		typeName, testCaseData,
	)
}

type TestDataKind uint8

const (
	In TestDataKind = iota
	Out
)

func (v *Helper) panicOnErrOfGetValues(err error, testeeName string, testDataKind TestDataKind, idx int) {
	bwerr.Panic(
		"<ansiPath>%s<ansi> expects testCase.%s[%d] "+err.Error(),
		testeeName, testDataKind.String(), idx,
	)
}

func (v *Helper) TestItemString(ti TestItem) (result string) {
	result = "_" + v.IdSet + "TestItem" + strings.ToUpper(ti.String())
	return result
}
