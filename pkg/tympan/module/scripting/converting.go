package scripting

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/FlagrantGarden/flfa/pkg/tympan/utils"
	"github.com/mitchellh/mapstructure"
)

// MetaConfig structs are used when parsing tags on configuration structs; they help turn a mapstructure tag into the
// name of a viper configuration key and change the behavior of a configuration item via the tympanconfig directive;
// right now the only supported directive is `tympanconfig:"ignore"` which ensures a struct key is not written to the
// configuration.
type MetaConfig struct {
	ConfigKey string
	Ignore    bool
}

// ParseStructTags() is used to introspect on a struct which is to be converted to a map Tengo can user; it returns the
// MetaConfig for a given struct field which SetStruct uses to determine behavior.
func ParseStructTags(tagEntry reflect.StructTag) (metaConfig MetaConfig) {
	mapstructTag, ok := tagEntry.Lookup("mapstructure")
	if ok {
		ignoreEntries := []string{"squash", "remain", "omitempty"}
		for _, entry := range strings.Split(mapstructTag, ",") {
			if !utils.Contains(ignoreEntries, entry) {
				metaConfig.ConfigKey = entry
			}
		}
	}
	flfaTag, ok := tagEntry.Lookup("flfa")
	if ok {
		tympanConfigDirectives := strings.Split(flfaTag, ",")
		if utils.Contains(tympanConfigDirectives, "ignore") {
			metaConfig.Ignore = true
		}
	}
	return
}

// ConvertToTengoMap is a helper function for intelligently transforming an arbitrary struct into a map[string]any
// object, which tengo can treat as a basic map. It is aware of both the mapstructure and tympanconfig directives.
func ConvertToTengoMap(data any) (map[string]any, error) {
	tengoMap := make(map[string]any)
	reflected_value := reflect.ValueOf(data)
	if reflected_value.Kind() == reflect.Pointer {
		reflected_value = reflect.Indirect(reflected_value)
	}
	reflected_type := reflected_value.Type()
	if reflected_value.Kind() != reflect.Struct {
		return tengoMap, fmt.Errorf("cannot set '%s' with SetStruct: expected a struct, got '%s'", reflected_type, reflected_value.Kind())
	}

	// Loop over the fields, adding any values to the map
	for _, field := range reflect.VisibleFields(reflected_type) {
		meta := ParseStructTags(field.Tag)

		// don't write this struct field to the map
		if meta.Ignore {
			continue
		}

		// Get the value, move from pointer to value if needed
		value := reflected_value.FieldByIndex(field.Index)
		if value.Kind() == reflect.Pointer {
			value = reflect.Indirect(value)
		}

		// ignore for now; skip invalid (zero value, no point in writing)
		if value.Kind() == reflect.Invalid || value.IsZero() {
			continue
		}

		// determine the key name for the map, replacing downcased field with value from mapstructure if specified
		name := field.Name
		if meta.ConfigKey != "" {
			name = meta.ConfigKey
		}
		name = strings.ToLower(name)

		// recurse if a struct, otherwise set the value
		if value.Kind() == reflect.Struct {
			mapValue, err := ConvertToTengoMap(value.Interface())
			if err != nil {
				return map[string]any{}, err
			}
			tengoMap[name] = mapValue
		} else if value.Kind() == reflect.Slice {
			// slices must be turned into []any
			array := make([]any, value.Len())
			for i := 0; i < value.Len(); i++ {
				array[i] = value.Index(i).Interface()
			}
			tengoMap[name] = array
		} else {
			tengoMap[name] = value.Interface()
		}
	}

	return tengoMap, nil
}

// ConvertFromTengoMap reverses the process, turning an arbitrary tengo map into a defined struct so that you can pass
// information back and forth between tengo and your application.
func ConvertFromTengoMap[T any](tengoMap map[string]any) (data T, err error) {
	err = mapstructure.Decode(tengoMap, &data)

	return data, err
}

// out of time but: should be returning a &tengo.Map with maps for child structs. We need to convert every single value
// to the appropriate tengo value. Probably should just have a `ToScriptValue()` and `FromScriptValue() method for
// every struct.
