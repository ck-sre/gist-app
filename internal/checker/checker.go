package checker

import (
	"slices"
	"strings"
	"unicode/utf8"
)

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
