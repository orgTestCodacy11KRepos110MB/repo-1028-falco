package test

import (
	"fmt"
	"regexp"

	"github.com/falcosecurity/falco/regression-tests/pkg/falco/run"
)

/*
addl_cmdline_opts
check_detection_counts
detect_counts
detect_level
disable_tags
disabled_rules
enable_source
json_include_output_property
json_include_tags_property
json_output
outputs
package
priority
rules_events
rules_file
run_duration
run_tags
should_detect
time_iso_8601
trace_file
validate_errors
validate_json
validate_ok
validate_rules_file
validate_warnings
*/

func Config(c *run.Config) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			ts.config = *c
		}
		return nil
	}
}

func Args(args ...string) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			ts.args = args
		}
		return nil
	}
}

func ExitCode(v int) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			return nil
		}
		if ts.err != nil {
			if codeErr, ok := ts.err.(*run.CodeError); ok {
				if codeErr.Code != v {
					return fmt.Errorf("expected exit code %d, but got %d", v, codeErr.Code)
				}
				return nil
			}
		}
		if v == 0 {
			return nil
		}
		return fmt.Errorf("expected exit code %d, but could not retrieve it", v)
	}
}

func Stdout(f func(string) error) Condition {
	return func(ts *testerState) error {
		if ts.done {
			return f(ts.stdout.String())
		}
		return nil
	}
}

func Stderr(f func(string) error) Condition {
	return func(ts *testerState) error {
		if ts.done {
			return f(ts.stderr.String())
		}
		return nil
	}
}

func StdoutMatch(e *regexp.Regexp) Condition {
	return Stdout(func(s string) error {
		if !e.MatchString(s) {
			return fmt.Errorf("stdout does not match regular expression:\nregexp=%s\n\nstdout=%s", e.String(), s)
		}
		return nil
	})
}

func StdoutNotMatch(e *regexp.Regexp) Condition {
	return Stdout(func(s string) error {
		if e.MatchString(s) {
			return fmt.Errorf("stdout does matches regular expression:\nregexp=%s\n\nstdout=%s", e.String(), s)
		}
		return nil
	})
}

func StderrMatch(e *regexp.Regexp) Condition {
	return Stderr(func(s string) error {
		if !e.MatchString(s) {
			return fmt.Errorf("stderr does not match regular expression:\nregexp=%s\n\nstderr=%s", e.String(), s)
		}
		return nil
	})
}

func StderrNotMatch(e *regexp.Regexp) Condition {
	return Stderr(func(s string) error {
		if e.MatchString(s) {
			return fmt.Errorf("stderr does matches regular expression:\nregexp=%s\n\nstderr=%s", e.String(), s)
		}
		return nil
	})
}
