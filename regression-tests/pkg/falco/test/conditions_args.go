package test

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

func AllEvents(v bool) Condition {
	return func(ts *testerState) error {
		if !ts.done {
			found := false
			var opts []string
			for _, o := range ts.args {
				if o != "-A" {
					opts = append(opts, o)
				} else {
					if v {
						opts = append(opts, o)
						found = true
					}
				}
			}
			if !found {
				opts = append(opts, "-A")
			}
			ts.args = opts
		}
		return nil
	}
}
