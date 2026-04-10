//go:build linux

package main

import (
	"net/http"
	"os"
	"time"

	daecmd "github.com/daeuniverse/dae/cmd"
	daejson "github.com/daeuniverse/dae/common/json"
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"

	"github.com/qi-mooo/dae-next/internal/daecontroller"
	"github.com/qi-mooo/dae-next/internal/uiassets"
)

var Version = "dev"

func main() {
	jsoniter.RegisterTypeDecoder("bool", &daejson.FuzzyBoolDecoder{})
	extra.RegisterFuzzyDecoders()

	http.DefaultClient.Timeout = 30 * time.Second
	uiassets.EnsureEnv()

	daecmd.Version = Version
	daecmd.ProgramName = "dae-next"
	daecmd.ProgramShort = "dae-next is dae packaged with the next-ui workspace."
	daecmd.ProgramLong = `dae-next is dae packaged with the next-ui workspace.`
	daecmd.ControllerFactory = daecontroller.New

	if err := daecmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
