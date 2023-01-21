package test

import (
	"encoding/json"
	"strings"
	"time"
)

type Alert struct {
	Hostname     string                 `json:"hostname"`
	Output       string                 `json:"output"`
	Priority     string                 `json:"priority"`
	Rule         string                 `json:"rule"`
	Source       string                 `json:"source"`
	Tags         []string               `json:"tags"`
	Time         time.Time              `json:"time"`
	OutputFields map[string]interface{} `json:"output_fields"`
}

func Detections(funcs ...func([]*Alert) error) Condition {
	return joinConditions(
		OutputJSON(true),
		editConfig("stdout_output:\n  enabled: true"),
		Stdout(func(stdout string) error {
			var alerts []*Alert
			err := readLineByLine(strings.NewReader(stdout), func(s string) error {
				var a Alert
				err := json.Unmarshal(([]byte)(s), &a)
				if err == nil {
					alerts = append(alerts, &a)
				}
				return nil // todo: find a way to not ignore non-json outputs
			})
			if err == nil {
				for _, f := range funcs {
					err = f(alerts)
					if err != nil {
						break
					}
				}
			}
			return err
		}),
	)
}

var AnyAlert = func(*Alert) error {
	return nil
}

func AlertRule(f func(string) error) func(*Alert) error {
	return func(a *Alert) error {
		return f(a.Rule)
	}
}

func AlertPriorty(f func(string) error) func(*Alert) error {
	return func(a *Alert) error {
		return f(a.Priority)
	}
}

func AlertSource(f func(string) error) func(*Alert) error {
	return func(a *Alert) error {
		return f(a.Source)
	}
}
