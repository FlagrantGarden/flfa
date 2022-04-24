package terminal

// A DynamicStyleList is a map of named lists of operations, each key representing a dynamic style.
type DynamicStyleList map[string][]Operation

// Add one or more operations to a dynamic style; if the style does not exist, create it.
func (list DynamicStyleList) Add(key string, operations ...Operation) {
	if len(list[key]) > 0 {
		list[key] = append(list[key], operations...)
	} else {
		list[key] = operations
	}
}

// Replace the definition of a dynamic style or create it if it does not exist.
func (list DynamicStyleList) Set(key string, operations ...Operation) {
	list[key] = operations
}

// Remove a dynamic style from the list by name.
func (list DynamicStyleList) Delete(key string) {
	delete(list, key)
}

// This terminal settings option sets a dynamic style for a new instance of terminal settings,
// creating it if it doesn't exist or overwriting it if it does.
func WithDynamicStyle(name string, operations ...Operation) Option {
	return func(settings *Settings) *Settings {
		settings.DynamicStyles.Set(name, operations...)
		return settings
	}
}
