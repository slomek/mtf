package core

import (
	"flag"
)

type ArgSettings struct {
	BuildBinary             bool
	StopComponentsAfterExit bool
}

var Settings = ArgSettings{}

func init() {
	flag.BoolVar(&Settings.BuildBinary, "rebuild_binary", false,
		"Determin if SUT binary should be rebuilded before start execution started")

	flag.BoolVar(&Settings.StopComponentsAfterExit, "stop_components", false,
		"Don't stop components after test execution have been finished")

}
