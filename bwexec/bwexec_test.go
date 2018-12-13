package bwexec_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/baza-winner/bwcore/bwexec"
	"github.com/baza-winner/bwcore/bwtesting"
	"github.com/baza-winner/bwcore/bwval"
)

func TestMain(m *testing.M) { // https://stackoverflow.com/questions/23729790/how-can-i-do-test-setup-using-the-testing-package-in-go/34102842#34102842
	mySetupFunction()
	retCode := m.Run()
	// myTeardownFunction()
	os.Exit(retCode)
}

// var installOpt = defparse.MustParseMap("{v: 'err', exitOnError: true, s: 'none'}")
// var installOpt map[string]interface{}

var installOpt, testOpt interface{}

func init() {
	installOpt = bwval.From(bwval.S{S: `{
		verbosity "err"
		exitOnError true
		silent "none"
	}`}).Val
	testOpt = bwval.From(bwval.S{S: `{
		captureStdout true
		captureStderr true
	}`}).Val
}

func mySetupFunction() {
	bwexec.MustCmd(bwexec.Args(`go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper`), installOpt)
	bwexec.MustCmd(bwexec.Args(`go`, `install`, `github.com/baza-winner/bwcore/bwexec/bwexectesthelper2`), installOpt)
}

func ExampleCmd() {
	ret := bwexec.MustCmd(bwexec.Args(`bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`))
	fmt.Printf("%s: %v\n", "Stdout", ret.Stdout)
	fmt.Printf("%s: %v\n", "Stderr", ret.Stderr)
	fmt.Printf("%s: %v\n", "Output", ret.Output)
	fmt.Printf("%s: %v\n", "ExitCode", ret.ExitCode)
	// Unordered ouput:
	// - Stdout: [some]
	// - Stderr: [thing]
	// - Output: [some thing]
	// - ExitCode: 2
}

func TestCmd(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwexec.MustCmd,
		func() map[string]bwtesting.Case {
			tests := map[string]bwtesting.Case{
				"test1": {
					In: []interface{}{
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `none`, `-s`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
						testOpt,
					},
					Out: []interface{}{
						bwexec.CmdResult{
							Stdout: []string{
								"===== exitCode: 2",
								"===== stdout:",
								"some",
								"===== stderr:",
								"thing",
							},
						},
						// map[string]interface{}{
						// 	"stdout": []string{
						// 		"===== exitCode: 2",
						// 		"===== stdout:",
						// 		"some",
						// 		"===== stderr:",
						// 		"thing",
						// 	},
						// 	"stderr":   []string{},
						// 	"exitCode": 0,
						// },
					},
				},
				"test2": {
					In: []interface{}{
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `none`, `-s`, `all`, `-d`, `-n`, `-e`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
						testOpt,
					},
					Out: []interface{}{
						bwexec.CmdResult{
							// Stdout:   nil,
							// Stderr:   nil,
							ExitCode: 2,
						},
						// map[string]interface{}{
						// 	"stdout":   []string{},
						// 	"stderr":   []string{},
						// 	"exitCode": 2,
						// },
					},
				},
				"test3": {
					In: []interface{}{
						bwexec.Args(
							`bwexectesthelper2`,
							`-v`, `all`, `-d`, `-n`, `bwexectesthelper`, `-exit`, `2`, `<stdout>some<stderr>thing`,
						),
						testOpt,
					},
					Out: []interface{}{
						bwexec.CmdResult{
							Stdout: []string{
								"bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething . . .",
								"some",
								"ERR: bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething",
								"===== exitCode: 2",
								"===== stdout:",
								"some",
								"===== stderr:",
								"thing",
							},
							Stderr: []string{
								"thing",
							},
							ExitCode: 0,
						},
						// map[string]interface{}{
						// 	"stdout": []string{
						// 		"bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething . . .",
						// 		"some",
						// 		"ERR: bwexectesthelper -exit 2 \u003cstdout\u003esome\u003cstderr\u003ething",
						// 		"===== exitCode: 2",
						// 		"===== stdout:",
						// 		"some",
						// 		"===== stderr:",
						// 		"thing",
						// 	},
						// 	"stderr": []string{
						// 		"thing",
						// 	},
						// 	"exitCode": 0,
						// },
					},
				},
			}
			return tests
		}(),
		// "test1",
	)
}
