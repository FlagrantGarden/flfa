package terminal

// Flags represent arbitrary behavior toggles for terminal settings. A Flag can be set to on, set to off, or unset. The
// value of a flagged setting can be checked to change behavior when dealing with the terminal.
type Flag *bool

var (
	on  = true
	off = false

	// The flag is set to be on
	FlagOn = Flag(&on)
	// The flag is set to be off
	FlagOff = Flag(&off)
	// The Flag is not explicitly set
	FlagUnset = Flag(nil)
)

// A terminal settings option to add a flag by name and value to the list of flags for that setting. If the flag already
// exists, this option will override it.
func WithFlag(name string, setting Flag) Option {
	return func(settings *Settings) {
		settings.Flags[name] = setting
	}
}

// A shorthand option for specifying that a flag should be set to on for an instance of terminal settings.
func WithFlagOn(name string) Option {
	return WithFlag(name, FlagOn)
}

// A shorthand option for specifying that a flag should be set to off for an instance of terminal settings.
func WithFlagOff(name string) Option {
	return WithFlag(name, FlagOff)
}

// A shorthand option for specifying that a flag should be registered as unset for an instance of terminal settings.
func WithFlagUnset(name string) Option {
	return WithFlag(name, FlagUnset)
}

// FlagFromBool converts a boolean value into a valid Flag, returning On if true and Off if false.
func FlagFromBool(value bool) (flag Flag) {
	if value {
		flag = FlagOn
	} else {
		flag = FlagOff
	}
	return flag
}

// Sets a flag on an existing instance of terminal settings.
func (settings *Settings) SetFlag(name string, setting Flag) {
	settings.Flags[name] = setting
}

// Helper function to set a flag as On for an existing instance of terminal settings.
func (settings *Settings) SetFlagOn(name string) {
	settings.SetFlag(name, FlagOn)
}

// Helper function to set a flag as Off for an existing instance of terminal settings.
func (settings *Settings) SetFlagOff(name string) {
	settings.SetFlag(name, FlagOff)
}

// Helper function to mark a flag as Unset for an existing instance of terminal settings.
func (settings *Settings) UnsetFlag(name string) {
	settings.SetFlag(name, FlagUnset)
}

// Remove a flag entirely from an existing instance of terminal settings.
func (settings *Settings) RemoveFlag(name string) {
	delete(settings.Flags, name)
}

// Retrieve the current state of a flag by name from an instance of terminal settings. If the flag cannot be found, it
// is treated as being unset.
func (settings *Settings) Flag(name string) Flag {
	if setting, ok := settings.Flags[name]; ok {
		return setting
	}
	return FlagUnset
}
