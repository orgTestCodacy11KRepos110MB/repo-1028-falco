package test

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func Prints[T any](w io.Writer) func(T) error {
	return func(s T) error {
		fmt.Fprintf(w, "%v", s)
		return nil
	}
}

func Equals[T comparable](v T) func(T) error {
	return func(s T) error {
		if s != v {
			return fmt.Errorf("value is not equal to expected value:\nexpected=%v\n\nactual=%v", v, s)
		}
		return nil
	}
}

func NotEquals[T comparable](values ...T) func(T) error {
	return func(s T) error {
		for _, v := range values {
			if s == v {
				return fmt.Errorf("value is equal to unwanted value:\nunwanted=%v\n\nactual=%v", v, s)
			}
		}
		return nil
	}
}

func joinChecks[T any](v T, funcs ...func(T) error) error {
	for _, f := range funcs {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func Count[T any](check func(T) error, funcs ...func(int) error) func([]T) error {
	return func(values []T) error {
		count := 0
		for _, a := range values {
			if err := check(a); err == nil {
				count++
			}
		}
		return joinChecks(count, funcs...)
	}
}

func Where[T any](check func(T) error, funcs ...func([]T) error) func([]T) error {
	return func(values []T) error {
		var subset []T
		for _, a := range values {
			if err := check(a); err == nil {
				subset = append(subset, a)
			}
		}
		return joinChecks(subset, funcs...)
	}
}

func Some[T any](check func(T) error, funcs ...func(T) error) func([]T) error {
	return func(values []T) error {
		var err error
		for _, a := range values {
			if err = check(a); err == nil {
				if err = joinChecks(a, funcs...); err == nil {
					return nil
				}
			}
		}
		return fmt.Errorf("no element matches the check, the last issues is: %s", err.Error())
	}
}

func Contains(values ...string) func(string) error {
	return func(s string) error {
		for _, v := range values {
			if !strings.Contains(s, v) {
				return fmt.Errorf("text does not contain subtext:\nsubtext=%s\n\ntext=%s", v, s)
			}
		}
		return nil
	}
}

func NotContains(values ...string) func(string) error {
	return func(s string) error {
		for _, v := range values {
			if strings.Contains(s, v) {
				return fmt.Errorf("text contains unwanted subtext:\nsubtext=%s\n\ntext=%s", v, s)
			}
		}
		return nil
	}
}

func Matches(e *regexp.Regexp) func(string) error {
	return func(s string) error {
		if !e.MatchString(s) {
			return fmt.Errorf("text does not match regular expression:\nregexp=%s\n\ntext=%s", e.String(), s)
		}
		return nil
	}
}

func NotMatches(exps ...*regexp.Regexp) func(string) error {
	return func(s string) error {
		for _, e := range exps {
			if e.MatchString(s) {
				return fmt.Errorf("text matches regular expression:\nregexp=%s\n\ntext=%s", e.String(), s)
			}
		}
		return nil
	}
}

func EncodesJSON(funcs ...func(map[string]interface{}) error) func(string) error {
	return func(s string) error {
		obj := make(map[string]interface{})
		err := json.Unmarshal(([]byte)(s), &obj)
		if err != nil {
			return fmt.Errorf("text is not in JSON format\ntext=%s", s)
		}
		for _, f := range funcs {
			if err = f(obj); err != nil {
				return err
			}
		}
		return nil
	}
}

// todo: support json pointer or something similar for nested keys
func Maps(key string, funcs ...func(string) error) func(map[string]interface{}) error {
	return func(m map[string]interface{}) error {
		if v, ok := m[key]; ok {
			if len(funcs) > 0 {
				s := fmt.Sprintf("%v", v)
				for _, f := range funcs {
					if err := f(s); err != nil {
						return err
					}
				}
			}
			return nil
		}
		return fmt.Errorf("map does not have key: %s", key)
	}
}
