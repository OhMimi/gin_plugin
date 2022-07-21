package service

import (
	"regexp"
	"strings"
)

var matchAllCap = regexp.MustCompile(`([a-z\d])([A-Z])`)

// CamelCaseToSnakeCase camel case convert to snake case
func CamelCaseToSnakeCase(str string) string {
	snake := matchAllCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}
