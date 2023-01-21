package test

import (
	"fmt"
	"regexp"
	"strings"
)

func Equal[T comparable](v T) func(T) error {
	return func(s T) error {
		if s != v {
			return fmt.Errorf("value is not equal to expected value:\nexpected=%v\n\nactual=%v", v, s)
		}
		return nil
	}
}

func NotEqual[T comparable](v T) func(T) error {
	return func(s T) error {
		if s == v {
			return fmt.Errorf("value is equal to unwanted value:\nunwanted=%v\n\nactual=%v", v, s)
		}
		return nil
	}
}

func Contain(v string) func(string) error {
	return func(s string) error {
		if !strings.Contains(s, v) {
			return fmt.Errorf("text does not contain subtext:\nsubtext=%s\n\ntext=%s", v, s)
		}
		return nil
	}
}

func NotContain(v string) func(string) error {
	return func(s string) error {
		if strings.Contains(s, v) {
			return fmt.Errorf("text contains unwanted subtext:\nsubtext=%s\n\ntext=%s", v, s)
		}
		return nil
	}
}

func Match(e *regexp.Regexp) func(string) error {
	return func(s string) error {
		if !e.MatchString(s) {
			return fmt.Errorf("text does not match regular expression:\nregexp=%s\n\ntext=%s", e.String(), s)
		}
		return nil
	}
}

func NotMatch(e *regexp.Regexp) func(string) error {
	return func(s string) error {
		if e.MatchString(s) {
			return fmt.Errorf("text matches regular expression:\nregexp=%s\n\ntext=%s", e.String(), s)
		}
		return nil
	}
}
