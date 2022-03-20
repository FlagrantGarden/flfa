package utils

import (
	"github.com/gosimple/slug"
)

// Defines the options that can be passed to ValidFileNameWithOptions()
type ValidFileNameOptions struct {
	// Language sets the language to use for slug creation. If unknown or unspecified, defaults to "en"
	Language string
	// CustomSub stores custom substitution map
	CustomSub *map[string]string
	// CustomRuneSub stores custom rune substitution map
	CustomRuneSub *map[rune]string

	// MaxLength stores maximum slug length.
	// It's smart so it will cat slug after full word.
	// By default slugs aren't shortened.
	// If MaxLength is smaller than length of the first word, then returned
	// slug will contain only substring from the first word truncated
	// after MaxLength.
	MaxLength int

	// Lowercase defines if the resulting slug is transformed to lowercase.
	// Default is true.
	Lowercase bool

	// Persistant defines whether the options should be reset after use of not.
	// If Persistant is true, the options will not be reset.
	Persistant bool
}

// Returns a valid utf-8 string, downcased and with special characters replaced or removed which might prevent the name
// from being a valid file name. It also replaces spaces with dashes.
func ValidFileName(name string) string {
	return slug.Make(name)
}

// Returns a valid utf-8 string, processed with the specified options. Useful in particular for when you want to limit
// the length of the filename, localize symbol substitutions, or prevent the file name from being downcased.
func ValidFileNameWithOptions(name string, options ValidFileNameOptions) (slugified string) {
	setSlugOptions(options)

	if options.Language != "" {
		slugified = slug.MakeLang(name, options.Language)
	} else {
		slugified = slug.Make(name)
	}

	if !options.Persistant {
		resetSlugOptions()
	}

	return slugified
}

// Returns a ValidFileNameOptions struct pre-populated with the defaults for the slug library so that the library's
// variables can be reset at the end of an invocation.
func getSlugDefaultOptions() ValidFileNameOptions {
	// I wish there was a better way to do this, but for now we just need to keep in sync with source:
	// https://pkg.go.dev/github.com/gosimple/slug#pkg-variables
	return ValidFileNameOptions{
		MaxLength: 0,
		Lowercase: true,
	}
}

// Resets the slug library's variables, which configure its behavior, to their defaults.
func resetSlugOptions() {
	setSlugOptions(getSlugDefaultOptions())
}

// Uses the specified options to set the slug library's variables to configure its behavior.
func setSlugOptions(options ValidFileNameOptions) {
	// If  a custom sub map is nil, set to empty.
	if options.CustomSub == nil {
		slug.CustomSub = map[string]string{}
	} else {
		slug.CustomSub = *options.CustomSub
	}
	if options.CustomRuneSub == nil {
		slug.CustomRuneSub = map[rune]string{}
	} else {
		slug.CustomRuneSub = *options.CustomRuneSub
	}
	slug.MaxLength = options.MaxLength
	slug.Lowercase = options.Lowercase
}
