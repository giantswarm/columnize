package columnize

import (
	"fmt"
	"strings"
)

type Config struct {
	// The string by which the lines of input will be split.
	Delim string

	// The string by which columns of output will be separated.
	Glue string

	// The string by which columns of output will be prefixed.
	Prefix string

	// A replacement string to replace empty fields
	Empty string
}

// Returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		Delim:  "|",
		Glue:   "  ",
		Prefix: "",
	}
}

// Returns a list of elements, each representing a single item which will
// belong to a column of output.
func getElementsFromLine(config *Config, line string) []interface{} {
	elements := make([]interface{}, 0)
	for _, field := range strings.Split(line, config.Delim) {
		value := strings.TrimSpace(field)
		if value == "" && config.Empty != "" {
			value = config.Empty
		}
		elements = append(elements, value)
	}
	return elements
}

// Examines a list of strings and determines how wide each column should be
// considering all of the elements that need to be printed within it.
func getWidthsFromLines(config *Config, lines []string) []int {
	var widths []int

	for _, line := range lines {
		elems := getElementsFromLine(config, line)
		for i := 0; i < len(elems); i++ {
			// remove color code for counting the length
			l := len(removeColorCode(elems[i].(string)))
			if len(widths) <= i {
				widths = append(widths, l)
			} else if widths[i] < l {
				widths[i] = l
			}
		}
	}
	return widths
}

// Given a set of column widths and the number of columns in the current line,
// returns a sprintf-style format string which can be used to print output
// aligned properly with other lines using the same widths set.
// func (c *Config) getStringFormat(widths []int, columns int) string {
// 	// Start with the prefix, if any was given.
// 	stringfmt := c.Prefix

// 	// Create the format string from the discovered widths
// 	for i := 0; i < columns && i < len(widths); i++ {
// 		if i == columns-1 {
// 			stringfmt += "%s\n"
// 		} else {
// 			stringfmt += fmt.Sprintf("%%-%ds%s", widths[i], c.Glue)
// 		}
// 	}
// 	return stringfmt
// }

func (c *Config) getStringFormat(widths []int, elems []interface{}) string {
	// Start with the prefix, if any was given.
	stringfmt := c.Prefix

	// Create the format string from the discovered widths
	for i := 0; i < len(elems) && i < len(widths); i++ {
		if i == len(elems)-1 {
			stringfmt += "%s\n"
		} else {
			if containsColorCode(elems[i]) {
				stringfmt += fmt.Sprintf("%%-%ds%s", widths[i]+colorCodeLen(elems[i].(string)), c.Glue)
			} else {
				stringfmt += fmt.Sprintf("%%-%ds%s", widths[i], c.Glue)
			}
		}
	}
	return stringfmt
}

// MergeConfig merges two config objects together and returns the resulting
// configuration. Values from the right take precedence over the left side.
func MergeConfig(a, b *Config) *Config {
	var result Config = *a

	// Return quickly if either side was nil
	if a == nil || b == nil {
		return &result
	}

	if b.Delim != "" {
		result.Delim = b.Delim
	}
	if b.Glue != "" {
		result.Glue = b.Glue
	}
	if b.Prefix != "" {
		result.Prefix = b.Prefix
	}
	if b.Empty != "" {
		result.Empty = b.Empty
	}

	return &result
}

// Format is the public-facing interface that takes either a plain string
// or a list of strings and returns nicely aligned output.
func Format(lines []string, config *Config) string {
	var result string

	conf := MergeConfig(DefaultConfig(), config)
	widths := getWidthsFromLines(conf, lines)

	// Create the formatted output using the format string
	for _, line := range lines {
		elems := getElementsFromLine(conf, line)
		stringfmt := conf.getStringFormat(widths, elems)
		result += fmt.Sprintf(stringfmt, elems...)
	}

	// Remove trailing newline without removing leading/trailing space
	if n := len(result); n > 0 && result[n-1] == '\n' {
		result = result[:n-1]
	}

	return result
}

// Convenience function for using Columnize as easy as possible.
func SimpleFormat(lines []string) string {
	return Format(lines, nil)
}

func containsColorCode(i interface{}) bool {
	s, ok := i.(string)
	if ok {
		return strings.Contains(s, "\x1b[")
	}
	return false
}

func colorCodeLen(s string) int {
	withoutColorCodes := removeColorCode(s)
	return len(s) - len(withoutColorCodes)
}

func removeColorCode(s string) string {
	if strings.Contains(s, "\x1b[") {
		start := strings.Index(s, "m") + 1
		end := strings.Index(s, "\x1b[0")
		return s[start:end]
	}
	return s
}
