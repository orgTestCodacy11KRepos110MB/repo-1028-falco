package test

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/yaml.v3"
)

/*
ALERTS:
	check_detection_counts
	detect_counts
	detect_level
	should_detect
RULE_VALIDATION:
	validate_errors
	validate_json
	validate_ok
	validate_rules_file
	validate_warnings
CLI:
	* addl_cmdline_opts
	* disable_tags
	* run_tags
	* disabled_rules
	* enable_source
	* run_duration
	* trace_file
	* rules_files (also CONFIG)
CONFIG:
	(skip) json_include_output_property
	(skip) json_include_tags_property
	* json_output
	* priority
	(skip) time_iso_8601
DEPRECATED:
	(skip) rules_events
*/

func joinConditions(conditions ...Condition) Condition {
	return func(ts *testerState) error {
		for _, c := range conditions {
			if err := c(ts); err != nil {
				return err
			}
		}
		return nil
	}
}

func mergeMaps(a, b map[interface{}]interface{}) {
	for k, v := range b {
		// If you use map[string]interface{}, ok is always false here.
		// Because yaml.Unmarshal will give you map[interface{}]interface{}.
		if v, ok := v.(map[interface{}]interface{}); ok {
			if bv, ok := a[k]; ok {
				if bv, ok := bv.(map[interface{}]interface{}); ok {
					mergeMaps(bv, v)
					continue
				}
			}
		}
		a[k] = v
	}
}

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
			if add && !found {
				newArgs = append(newArgs, args...)
			}
			ts.args = newArgs
		}
		return nil
	}
}

func toggleTags(enable bool, tags ...string) Condition {
	arg := "-t"
	if !enable {
		arg = "-T"
	}
	var conds []Condition
	for _, t := range tags {
		conds = append(conds, addOrRemoveArg(true, false, arg, t))
	}
	return joinConditions(conds...)
}

func editConfig(config string) Condition {
	return func(ts *testerState) (err error) {
		if ts.done {
			return nil
		}
		curConfig := make(map[interface{}]interface{})
		editConfig := make(map[interface{}]interface{})
		if err = yaml.Unmarshal([]byte(ts.config), &curConfig); err == nil {
			if err = yaml.Unmarshal([]byte(config), &editConfig); err == nil {
				var newConf []byte
				mergeMaps(curConfig, editConfig)
				newConf, err = yaml.Marshal(&curConfig)
				if err == nil {
					ts.config = string(newConf)
				}
			}
		}
		return err
	}
}

func AllEvents(v bool) Condition {
	return addOrRemoveArg(v, true, "-A")
}

func EnableTags(tags ...string) Condition {
	return toggleTags(true, tags...)
}

func DisableTags(tags ...string) Condition {
	return toggleTags(true, tags...)
}

func DisableRules(rules ...string) Condition {
	var conds []Condition
	for _, r := range rules {
		conds = append(conds, addOrRemoveArg(true, false, "-D", r))
	}
	return joinConditions(conds...)
}

func EnableSources(sources ...string) Condition {
	var conds []Condition
	for _, s := range sources {
		conds = append(conds, addOrRemoveArg(true, false, "--enable-source", s))
	}
	return joinConditions(conds...)
}

func DisableSources(sources ...string) Condition {
	var conds []Condition
	for _, s := range sources {
		conds = append(conds, addOrRemoveArg(true, false, "--disable-source", s))
	}
	return joinConditions(conds...)
}

func RuleFiles(paths ...string) Condition {
	var conds []Condition
	for _, s := range paths {
		conds = append(conds, addOrRemoveArg(true, false, "-r", s))
	}
	return joinConditions(conds...)
}

func RulePriority(p string) Condition {
	return editConfig("priority: " + p)
}

func OutputJSON(v bool) Condition {
	if v {
		return editConfig("json_output: true")
	}
	return editConfig("json_output: false")
}

func CaptureFile(path string) Condition {
	return addOrRemoveArg(true, true, "-e", path)
}

func Duration(d time.Duration) Condition {
	return joinConditions(
		addOrRemoveArg(true, true, "-M", fmt.Sprintf("%d", int64(d.Seconds()))),
		func(ts *testerState) error {
			if !ts.done {
				ts.duration = skewedDuration(d)
			} else {
				if ts.err == context.DeadlineExceeded {
					return fmt.Errorf("falco did not terminate within expected duration: %s", d.String())
				}
			}
			return nil
		},
	)
}
