package bwos_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/baza-winner/bwcore/ansi"
	"github.com/baza-winner/bwcore/bw"
	"github.com/baza-winner/bwcore/bwos"
	"github.com/baza-winner/bwcore/bwtesting"
)

func ExampleShortenFileSpec() {
	fmt.Printf(`%q`, bwos.ShortenFileSpec(os.Getenv(`HOME`)+`/bw`))
	// Output: "~/bw"
}

func ExampleShortenFileSpec_2() {
	fmt.Printf(`%q`, bwos.ShortenFileSpec(`/lib/bw`))
	// Output: "/lib/bw"
}

func TestExitMsg(t *testing.T) {
	bwtesting.BwRunTests(t,
		bwos.ExitMsg,
		map[string]bwtesting.Case{
			"with newline": {
				In: []interface{}{
					bw.A{Fmt: "some exit msg\n"},
				},
				Out: []interface{}{
					ansi.String("some exit msg\n"),
				},
			},
			"without newline": {
				In: []interface{}{
					bw.A{Fmt: "some exit msg"},
				},
				Out: []interface{}{
					ansi.String("some exit msg") + "\n",
				},
			},
		},
	)
}
