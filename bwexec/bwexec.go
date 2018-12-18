// Package bwexec предоставялет функцию ExecCmd.
package bwexec

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bwerr"
	"github.com/baza-winner/bwcore/bwrune"
	"github.com/baza-winner/bwcore/bwstr"
	"github.com/baza-winner/bwcore/bwval"
)

const defaultFailedCode = 1

var cmdOptDef *bwval.Def

func init() {
	cmdOptDef = bwval.MustDefFrom(bwrune.S{S: `
		{
			type Map
			keys {
				verbosity {
					type String
					enum <all err ok none>
					default none
				}
				captureStdout {
					type Bool
					default false
				}
				captureStderr {
					type Bool
					default false
				}
				captureOutput {
					type Bool
					default false
				}
				silent {
					type String
					enum <none stderr stdout all>
					default all
				}
				exitOnError {
					type Bool
					default false
				}
				workDir {
					type String
					isOptional true
				}
			}
		}
	`})
}

func CmdOptDef() *bwval.Def {
	return cmdOptDef
}

// ============================================================================

type A struct {
	Cmd  string
	Args []string
}

// ============================================================================
func Args(cmdName string, cmdArgs ...string) A {
	return A{Cmd: cmdName, Args: cmdArgs}
}

// ============================================================================

type CmdResult struct {
	ExitCode int
	Stdout   []string
	Stderr   []string
	Output   []string
}

// ============================================================================

type CmdOpt struct {
	h bwval.Holder
}

func MustCmdOpt(a bwval.FromProvider) CmdOpt {
	return CmdOpt{bwval.MustFrom(a, bwval.O{
		Def:          cmdOptDef,
		PathProvider: bwval.PathS{S: "CmdOpt"},
	})}
}

func cmdOpt(optOpt []CmdOpt) (result CmdOpt) {
	if len(optOpt) > 0 {
		result = optOpt[0]
	} else {
		result = CmdOpt{bwval.MustFrom(
			bwval.V{},
			bwval.O{
				Def:          cmdOptDef,
				PathProvider: bwval.PathS{S: "CmdOpt"},
			},
		)}
	}
	return
}

func Cmd(a A, optOpt ...CmdOpt) (result CmdResult, err error) {
	opt := cmdOpt(optOpt)
	// var optFromProvider bwval.FromProvider
	// var opt.h CmdOpt
	// if len(optOpt) > 0 {
	// } else {

	// }

	// var opt bwval.FromProvider
	// if len(optOpt) > 0 {
	// 	opt = optOpt[0]
	// 	// if s, ok := opt.(string) {
	// 	// 	opt = bwval.
	// 	// } else {
	// 	// 	optFromProvider
	// 	// }
	// }
	// if opt == nil {
	// 	opt = bwval.V{Val: nil}
	// }
	// // var opt interface{}
	// // opt.h := bwval.MustFrom(bwval.V{Val: opt, Def: cmdOptDef}, bwval.PathS{S: "Cmd.opt"})
	// opt.h := bwval.MustFrom(opt, bwval.CmdOpt{
	// 	Def:          cmdOptDef,
	// 	PathProvider: bwval.PathS{S: "Cmd.opt"},
	// })

	cmdTitle := bwstr.SmartQuote(append([]string{a.Cmd}, a.Args...)...)
	optSilent := opt.h.MustPathStr("silent").MustString()
	optCaptureStdout := opt.h.MustPathStr("captureStdout").MustBool()
	optCaptureStderr := opt.h.MustPathStr("captureStderr").MustBool()
	optCaptureOutput := opt.h.MustPathStr("captureOutput").MustBool()
	optVerbosity := opt.h.MustPathStr("verbosity").MustString()
	optWorkDir := opt.h.MustPathStr("workDir?").MustString()
	var pwd string
	if optWorkDir != "" {
		if pwd, err = os.Getwd(); err != nil {
			return
		}
		if err = os.Chdir(optWorkDir); err != nil {
			return
		}
	}
	if optVerbosity == `all` || optVerbosity == `allBrief` {
		fmt.Println(ansi.String(`<ansiPath>` + cmdTitle + `<ansi> . . .`))
	}

	cmd := exec.Command(a.Cmd, a.Args...)

	type pipeStruct struct {
		getPipe           func() (io.ReadCloser, error)
		optCapturePipe    bool
		passPipe          bool
		pipeCaptureTarget *[]string
		pipe              *os.File
	}

	processPipe := func(a pipeStruct) (err error) {
		if optCaptureOutput || a.optCapturePipe || a.passPipe {
			var pipe io.ReadCloser
			if pipe, err = a.getPipe(); err != nil {
				return
			}
			scanner := bufio.NewScanner(pipe)
			go func() {
				for scanner.Scan() {
					s := scanner.Text()
					if a.optCapturePipe {
						*a.pipeCaptureTarget = append(*a.pipeCaptureTarget, s)
					}
					if optCaptureOutput {
						result.Output = append(result.Output, s)
					}
					if a.passPipe {
						fmt.Fprintln(a.pipe, s)
					}
				}
			}()
		}
		return
	}

	processPipe(pipeStruct{
		getPipe:           cmd.StdoutPipe,
		optCapturePipe:    optCaptureStdout,
		passPipe:          !(optSilent == `all` || optSilent == `stdout`),
		pipeCaptureTarget: &result.Stdout,
		pipe:              os.Stdout,
	})

	processPipe(pipeStruct{
		getPipe:           cmd.StderrPipe,
		optCapturePipe:    optCaptureStderr,
		passPipe:          !(optSilent == "all" || optSilent == "stderr"),
		pipeCaptureTarget: &result.Stderr,
		pipe:              os.Stderr,
	})

	if err = cmd.Start(); err != nil {
		return
	}

	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	if err = cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); !ok {
			return
		} else if status, ok := exiterr.Sys().(syscall.WaitStatus); !ok {
			return
			// log.Printf(ansi.String("<ansiWarn>Could not get exit code for failed program: <ansiPath>%s"), cmdTitle)
			// exitCode = defaultFailedCode
		} else {
			result.ExitCode = status.ExitStatus()
			err = nil
		}
	}

	var ansiName, prefix string
	if result.ExitCode == 0 && (optVerbosity == `all` || optVerbosity == `allBrief` || optVerbosity == `ok`) {
		ansiName, prefix = `ansiOK`, `OK`
	}
	if result.ExitCode != 0 && (optVerbosity == `all` || optVerbosity == `allBrief` || optVerbosity == `err`) {
		ansiName, prefix = `ansiErr`, `ERR`
	}
	if len(prefix) > 0 {
		fmt.Println(ansi.StringA(ansi.A{Default: ansi.MustTag(ansiName), S: prefix + `: <ansiPath>` + cmdTitle}))
	}
	if opt.h.MustPathStr("exitOnError").MustBool() && result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}
	if pwd != "" {
		if err = os.Chdir(pwd); err != nil {
			return
		}
	}
	return
}

// ============================================================================

func MustCmd(a A, optOpt ...CmdOpt) (result CmdResult) {
	var err error
	if result, err = Cmd(a, optOpt...); err != nil {
		bwerr.PanicErr(err)
	}
	return
}

// ============================================================================
