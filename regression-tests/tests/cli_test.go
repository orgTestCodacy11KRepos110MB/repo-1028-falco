package tests

import (
	"os/user"
	"regexp"
	"testing"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
	"github.com/falcosecurity/falco/regression-tests/pkg/falco/test"
)

var DefaultFalcoBinaryPath = "/usr/bin/falco"
var IsRoot = false

func init() {
	user, err := user.Current()
	if err != nil {
		panic("could not get the current user")
	}
	IsRoot = user.Username == "root"
}

func TestBadConfigNoOutput(t *testing.T) {
	test.RunTest(t,
		run.NewProcRunner(DefaultFalcoBinaryPath),
		test.ExitCode(1),
		test.StderrMatch(regexp.MustCompile(".*No outputs configured.*")),
	)
}

func TestArgVersion(t *testing.T) {
	c := run.Config{}
	c.StdoutOutput.Enabled = true
	test.RunTest(t,
		run.NewProcRunner(DefaultFalcoBinaryPath),
		test.Config(&c),
		test.Args("--version"),
		test.StdoutMatch(regexp.MustCompile(
			"Falco version:[\\s]+[0-9]+\\.[0-9]+\\.[0-9][\\s]+"+
				"Libs version:[\\s]+[0-9]+\\.[0-9]+\\.[0-9][\\s]+"+
				"Plugin API:[\\s]+[0-9]+\\.[0-9]+\\.[0-9][\\s]+"+
				"Driver:[\\s]+"+
				"API version:[\\s]+[0-9]+\\.[0-9]+\\.[0-9][\\s]+"+
				"Schema version:[\\s]+[0-9]+\\.[0-9]+\\.[0-9][\\s]+"+
				"Default driver:[\\s]+[0-9]+\\.[0-9]+\\.[0-9]\\+driver")),
		test.ExitCode(0),
	)
	// todo: test json output too
}
