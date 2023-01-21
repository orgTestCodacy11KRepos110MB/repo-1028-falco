package tests

import (
	"regexp"
	"testing"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
	"github.com/falcosecurity/falco/regression-tests/pkg/falco/test"
)

func TestFailEmptyConfig(t *testing.T) {
	test.RunTest(t, run.NewExecutableRunner(FalcoExecutable),
		test.ExitCode(test.Equals(1)),
		test.Stderr(test.Contains("You must specify at least one rules")),
	)
}

func TestCmdVersion(t *testing.T) {
	runner := run.NewExecutableRunner(FalcoExecutable)
	t.Run("text", func(t *testing.T) {
		test.RunTest(t, runner,
			test.Args("--version"),
			test.ExitCode(test.Equals(0)),
			test.Stdout(test.Matches(regexp.MustCompile(
				`Falco version:[\s]+[0-9]+\.[0-9]+\.[0-9](\-[0-9]+\+[a-f0-9]+)?[\s]+`+
					`Libs version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
					`Plugin API:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
					`Engine:[\s]+[0-9]+[\s]+`+
					`Driver:[\s]+`+
					`API version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
					`Schema version:[\s]+[0-9]+\.[0-9]+\.[0-9][\s]+`+
					`Default driver:[\s]+[0-9]+\.[0-9]+\.[0-9]\+driver`))),
		)
	})
	t.Run("json", func(t *testing.T) {
		test.RunTest(t, runner,
			test.Args("--version"),
			test.OutputJSON(true),
			test.ExitCode(test.Equals(0)),
			test.Stdout(
				test.EncodesJSON(
					test.Maps("default_driver_version"),
					test.Maps("driver_api_version"),
					test.Maps("driver_schema_version"),
					test.Maps("engine_version"),
					test.Maps("falco_version"),
					test.Maps(
						"libs_version",
						test.Equals("0.10.0"),
					),
					test.Maps("plugin_api_version"),
				),
			),
		)
	})
}
