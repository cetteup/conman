// Data structure for storing data from Refractong engine configuration files (.con)
package config

import (
	"fmt"
	"sort"
	"strings"
)

const (
	quoteChar           = "\""
	multiValueSeparator = ";"
)

type ErrNoSuchKey struct {
	path string
	key  string
}

func (e *ErrNoSuchKey) Error() string {
	return fmt.Sprintf("no such key in %s: %q", e.path, e.key)
}

type Config struct {
	Path    string
	content map[string]Value
}

func New(path string, content map[string]Value) *Config {
	return &Config{
		Path:    path,
		content: content,
	}
}

func FromBytes(path string, data []byte) *Config {
	// Split on \n in order to make parsing work with either \r\n or just \n line breaks
	lines := strings.Split(string(data), "\n")

	parsed := map[string]Value{}
	for _, line := range lines {
		// Trim any \r from line and split on first space
		elements := strings.SplitN(strings.Trim(line, "\r"), " ", 2)

		// TODO do something other than ignoring any invalid lines here?
		if len(elements) == 2 {
			// Add key, value or append to value
			key, content := elements[0], elements[1]
			current, exists := parsed[key]
			if exists {
				content = strings.Join([]string{current.content, content}, multiValueSeparator)
			}
			parsed[key] = *NewValue(content)
		}
	}

	return &Config{
		Path:    path,
		content: parsed,
	}
}

func (c *Config) HasKey(key string) bool {
	_, ok := c.content[key]
	return ok
}

func (c *Config) GetValue(key string) (Value, error) {
	value, ok := c.content[key]
	if !ok {
		return Value{}, &ErrNoSuchKey{
			path: c.Path,
			key:  key,
		}
	}
	return value, nil
}

func (c *Config) SetValue(key string, value Value) {
	c.content[key] = value
}

func (c *Config) Delete(key string) {
	delete(c.content, key)
}

func (c *Config) ToBytes() []byte {
	lines := make([]string, 0)

	for key, value := range c.content {
		// value.Slice() returns a single element slice for non-multi values, so we can safely iterate the slice even for those
		for _, subValue := range value.slice() {
			lines = append(lines, fmt.Sprintf("%s %s", key, subValue))
		}
	}

	// map iteration order is pseudo-random, so sort lines alphabetically to ensure we always generate the same byte array for a given config
	sort.Slice(lines, func(i, j int) bool {
		return strings.Compare(lines[i], lines[j]) < 1
	})

	// append an empty line, else BF2 will reset the configured value of the last line to default
	lines = append(lines, "")

	return []byte(strings.Join(lines, "\r\n"))
}

type Value struct {
	content string
}

func NewValue(content string) *Value {
	return &Value{
		content: content,
	}
}

func NewQuotedValue(content string) *Value {
	return NewValue(quoteValue(content))
}

func NewValueFromSlice(content []string) *Value {
	return &Value{
		content: strings.Join(content, multiValueSeparator),
	}
}

func NewQuotedValueFromSlice(content []string) *Value {
	quoted := make([]string, 0, len(content))
	for _, c := range content {
		quoted = append(quoted, quoteValue(c))
	}

	return NewValueFromSlice(quoted)
}

func (v *Value) String() string {
	if isQuotedValue(v.content) {
		return strings.Trim(v.content, quoteChar)
	}
	return v.content
}

func (v *Value) slice() []string {
	return strings.Split(v.content, multiValueSeparator)
}

func (v *Value) Slice() []string {
	values := v.slice()
	for i, item := range values {
		if isQuotedValue(item) {
			values[i] = strings.Trim(item, quoteChar)
		}
	}
	return values
}

// isQuotedValue Checks whether a config value is a quoted string (starts and ends with a quote character, with no other quote characters in between)
func isQuotedValue(value string) bool {
	return strings.HasPrefix(value, quoteChar) && strings.HasSuffix(value, quoteChar) && strings.Count(value, quoteChar) == 2
}

func quoteValue(value string) string {
	return fmt.Sprintf("%[1]s%s%[1]s", quoteChar, value)
}
