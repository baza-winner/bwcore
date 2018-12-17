// Package bwos содержит bw-дополнение для package os.
package bwos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwerr"
)

// ============================================================================

// ShortenFileSpec укорачивает строку за счет замены префикса, совпадающиего (если) cо значением
// ${HOME} (переменная среды), на символ `~`
func ShortenFileSpec(fileSpec string) string {
	fileSpec = filepath.Clean(fileSpec)
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		homeDir = AppendPathSeparator(filepath.Clean(homeDir))
		if len(fileSpec) >= len(homeDir) && fileSpec[:len(homeDir)] == homeDir {
			fileSpec = "~/" + fileSpec[len(homeDir):]
		}
	}
	return fileSpec
}

func IsInPath(dirSpec, fileSpec string) (ok bool, err error) {
	if dirSpec, err = filepath.Abs(dirSpec); err != nil {
		return
	}
	if fileSpec, err = filepath.Abs(fileSpec); err != nil {
		return
	}
	dirSpec = AppendPathSeparator(dirSpec)
	fileSpec = AppendPathSeparator(fileSpec)
	ok = len(fileSpec) >= len(dirSpec) && fileSpec[:len(dirSpec)] == dirSpec
	return
}

func AppendPathSeparator(dirSpec string) string {
	if len(dirSpec) > 0 && dirSpec[len(dirSpec)-1] != os.PathSeparator {
		dirSpec += string(os.PathSeparator)
	}
	return dirSpec
}

func Exit(exitCode int, fmtString string, fmtArgs ...interface{}) {
	ExitA(exitCode, bw.A{fmtString, fmtArgs})
}

// Exit with exitCode and message defined by bw.I
func ExitA(exitCode int, a bw.I) {
	fmt.Print(ExitMsg(a))
	os.Exit(exitCode)
}

// ============================================================================

var newlineAtTheEnd, _ = regexp.Compile(`\n\s*$`)

func ExitMsg(a bw.I) (result string) {
	err := bwerr.FromA(a)
	result = err.S
	if !newlineAtTheEnd.MatchString(ansi.ChopReset(result)) {
		result += string('\n')
	}
	return
}

// ============================================================================

func ResolveSymlink(fileSpec string, optDepth ...uint) (result string, err error) {
	var depth uint
	if len(optDepth) > 0 {
		depth = optDepth[0]
	}
	for {
		var fi os.FileInfo
		if fi, err = os.Lstat(fileSpec); err != nil {
			return
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			break
		} else if fileSpec, err = os.Readlink(fileSpec); err != nil {
			return
		} else if depth == 1 {
			break
		} else if depth > 1 {
			depth--
		}
	}
	result = fileSpec
	return
}

func IsSymlink(fileSpec string) (ok bool, err error) {
	var fi os.FileInfo
	if fi, err = os.Lstat(fileSpec); err == nil {
		ok = fi.Mode()&os.ModeSymlink != 0
	}
	return
}

func Exists(fileSpec string) (ok bool, err error) {
	if _, err = os.Stat(fileSpec); err == nil {
		ok = true
	} else if os.IsNotExist(err) {
		ok = false
		err = nil
	}
	return
}

func VerDir() (result string, err error) {
	var executableFileSpec string
	if executableFileSpec, err = os.Executable(); err != nil {
		return
	}
	linkSourceFileSpec := executableFileSpec
	for {
		var fi os.FileInfo
		if fi, err = os.Lstat(linkSourceFileSpec); err != nil {
			return
		}
		if fi.Mode()&os.ModeSymlink == 0 {
			break
		} else if linkSourceFileSpec, err = os.Readlink(linkSourceFileSpec); err != nil {
			return
		}
	}
	// bwdebug.Print("linkSourceFileSpec", linkSourceFileSpec)
	result = filepath.Clean(filepath.Join(filepath.Dir(linkSourceFileSpec), "..", ".."))
	return
}

// ============================================================================
