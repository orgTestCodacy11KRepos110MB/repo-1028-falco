package test

import "fmt"

/*
ALERTS:
	check_detection_counts
	detect_counts
	detect_level
	should_detect
CLI:
	addl_cmdline_opts
	disable_tags
	run_tags
	disabled_rules
	enable_source
	run_duration
	trace_file
	rules_files (also CONFIG)
CONFIG:
	json_include_output_property
	json_include_tags_property
	json_output
	priority
	time_iso_8601
RULE_VALIDATION:
	validate_errors
	validate_json
	validate_ok
	validate_rules_file
	validate_warnings
DEPRECATED:
	rules_events
*/

// addOrRemoveArg adds or remove an argument from the CLI (along with the
// other following arguments representing parameters, such as `-r <rule_filename>`.
// If add is false, the argument is removed from the CLI. Otherwise, the
// argument appended to the CLI if unique is false, and substitutes any
// previously-existing argument of the same kind if unique is true.
func addOrRemoveArg(add, unique bool, args ...string) Condition {
	return func(ts *testerState) error {
		if len(args) == 0 {
			return fmt.Errorf("bad test: tried adding or removing zero args")
		}
		if !ts.done {
			found := false
			var newArgs []string
			for i := 0; i < len(ts.args); i++ {
				arg := ts.args[i]
				if arg != args[0] {
					newArgs = append(newArgs, arg)
				} else {
					// the arg we want to add/remove is already present,
					// and may have some extra args following it
					if add {
						// check for repetition in case of uniqueness constraint
						if !unique || !found {
							// we want to substitute the existing arg with ours
							newArgs = append(newArgs, args...)
							found = true
						}
					}
					// we want to remove the existing arg from the CLI, so
					// we skip it and its following args
					i += len(args) - 1
				}
			}
			if !found {
				newArgs = append(newArgs, args...)
			}
			ts.args = newArgs
		}
		return nil
	}
}

func AllEvents(v bool) Condition {
	return addOrRemoveArg(v, true, "-A")
}
