package json

import (
	"encoding/json"
	"strings"

	"github.com/knadh/koanf/maps"
	"github.com/mitchellh/mapstructure"
)

// Configuration for how a given object should be converted to JSON
type Settings struct {
	// If specified, prefixes each line after the first with the value
	Prefix string
	// If specified, the JSON will be rendered indented, repeating this value once for each level of indent
	Indent string
	// A list of fields to ignore. You can specify dotpath notation, like "foo.bar"
	Ignore []string
}

// Options are functions you can pass to StructToJson to apply the settings dynamically
type Option func(settings *Settings)

// Sets the Prefix string when calling StructToJson
func WithPrefix(prefix string) Option {
	return func(settings *Settings) {
		settings.Prefix = prefix
	}
}

// Sets the Indent string when calling StructToJson
func WithIndent(indent string) Option {
	return func(settings *Settings) {
		settings.Indent = indent
	}
}

// Adds one or more dotpaths (or top-level keys) to be ignored when calling StructToJson
func WithIgnore(dotpaths ...string) Option {
	return func(settings *Settings) {
		settings.Ignore = append(settings.Ignore, dotpaths...)
	}
}

// A convenience function for printing an arbitrary struct as a JSON blob with/without indent, prefix, and dynamically
// ignored keys. Particularly useful when dealing with user-generated or otherwise dynamic/unpredictable structures.
func StructToJson(input any, options ...Option) string {
	var transitional map[string]any
	settings := &Settings{}

	for _, option := range options {
		option(settings)
	}

	mapstructure.Decode(input, &transitional)

	for _, dotpath := range settings.Ignore {
		maps.Delete(transitional, strings.Split(dotpath, "."))
	}

	output, _ := json.MarshalIndent(transitional, settings.Prefix, settings.Prefix)
	return string(output)
}
