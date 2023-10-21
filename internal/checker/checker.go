package checker

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Checker struct {
	AttrErrors map[string]string
}

func (c *Checker) CheckPassed() bool {
	return len(c.AttrErrors) == 0
}

func (c *Checker) AddAttrError(attr, msg string) {
	if c.AttrErrors == nil {
		c.AttrErrors = make(map[string]string)
	}

	if _, ok := c.AttrErrors[attr]; !ok {
		c.AttrErrors[attr] = msg
	}
}

func (c *Checker) CheckAttr(valid bool, attr, msg string) {
	if !valid {
		c.AddAttrError(attr, msg)
	}
}

func NotEmpty(val string) bool {
	return strings.TrimSpace(val) != ""
}

func LimitChars(val string, limit int) bool {
	return utf8.RuneCountInString(val) <= limit
}

func AllowedVal[U comparable](val U, allowedVals ...U) bool {
	return slices.Contains(allowedVals, val)
}

func CharMin(val string, min int) bool {
	return utf8.RuneCountInString(val) >= min
}

func StringMatches(val string, regex *regexp.Regexp) bool {
	return regex.MatchString(val)
}
