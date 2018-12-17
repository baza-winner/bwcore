// Вспомогательная утилита для тестирования bwexec.ExecCmd.
package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/baza-winner/bwcore/ansi"
	_ "github.com/baza-winner/bwcore/ansi/tags"
	"github.com/baza-winner/bwcore/bwexec"
	"github.com/baza-winner/bwcore/bwos"
	"github.com/baza-winner/bwcore/bwval"
)

func init() {
	ansi.MustAddTag("ansiHeader",
		ansi.MustSGRCodeOfColor8(ansi.Color8{Color: ansi.SGRColorWhite}),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdBold),
		ansi.MustSGRCodeOfCmd(ansi.SGRCmdFaint),
	)
}

func main() {
	exitOnErrorFlag := flag.Bool("e", false, "exitOnError")
	verbosityFlag := flag.String("v", "none", "verbosity: none, err, ok, all, allBrief")
	silentFlag := flag.String("s", "none", "silent: none, stderr, stdout, all")
	displayFlag := flag.Bool("d", false, "display exitCode, stdout, stderr, output")
	noColorFlag := flag.Bool("n", false, "no color")
	flag.Parse()
	if flag.NArg() < 1 {
		bwos.Exit(1, `<ansiPath>bwexectesthelper2<ansi> expects at least one arg`)
	}
	argsWithoutProg := flag.Args()

	display := *displayFlag
	ansi.NoColor = *noColorFlag

	ret := bwexec.MustCmd(
		bwexec.Args(argsWithoutProg[0], argsWithoutProg[1:]...),
		bwexec.MustCmdOpt(
			bwval.V{
				Val: map[string]interface{}{
					"verbosity":     *verbosityFlag,
					"exitOnError":   *exitOnErrorFlag,
					"silent":        *silentFlag,
					"captureStdout": true,
					"captureStderr": true,
				},
			},
		),
	)
	if display {
		ansiExitCode := `<ansiOK>`
		if ret.ExitCode != 0 {
			ansiExitCode = `<ansiErr>`
		}
		fmt.Printf(ansi.String(`<ansiHeader>===== exitCode: `+ansiExitCode+"%d\n"), ret.ExitCode)
		fmt.Println(ansi.String(`<ansiHeader>===== stdout:`))
		fmt.Println(strings.Join(ret.Stdout, "\n"))
		fmt.Println(ansi.String(`<ansiHeader>===== stderr:`))
		fmt.Println(strings.Join(ret.Stderr, "\n"))

		if len(ret.Output) > 0 {
			fmt.Println(ansi.String(`<ansiHeader>===== output:`))
			fmt.Println(strings.Join(ret.Output, "\n"))
		}
	}
}
