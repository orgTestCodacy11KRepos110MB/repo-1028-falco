package tests

import (
	"regexp"
	"testing"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
	"github.com/falcosecurity/falco/regression-tests/pkg/falco/test"
)

func TestBadConfigNoOutput(t *testing.T) {
	test.RunTest(t,
		run.NewProcRunner(DefaultFalcoExecutable),
		test.ExitCode(test.Equal(1)),
		test.Stderr(test.Contain("No outputs configured")),
	)
}

func TestArgVersion(t *testing.T) {
	test.RunTest(t,
		run.NewProcRunner(DefaultFalcoExecutable),
		test.Args("--version"),
		test.ExitCode(test.Equal(0)),
		test.Stdout(test.Match(regexp.MustCompile(
			`Falco version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
				`Libs version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
				`Plugin API:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
				`Driver:[\s]+`+
				`API version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
				`Schema version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
				`Default driver:[\s]+[0-9]+\.[0-9]+\.[0-9]\+driver`))),
	)
	// todo: test json output too
}
