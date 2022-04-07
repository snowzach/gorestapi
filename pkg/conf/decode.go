package conf

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
)

// DecodeOption is a helper function to configure a Decoder
type DecodeOption func(dc *mapstructure.DecoderConfig)

// Decode is a helper for github.com/mitchellh/mapstructure to compose
// a DecoderConfig decode into a struct.
func Decode(source, dest interface{}, opts ...DecodeOption) error {

	// Get the default and apply the options
	dc := DefaultDecoderConfig()
	for _, o := range opts {
		o(dc)
	}
	dc.Result = dest

	// Create the decoder and deccode into the destination.
	decoder, err := mapstructure.NewDecoder(dc)
	if err != nil {
		return fmt.Errorf("could not create decoder: %w", err)
	}
	return decoder.Decode(source)

}

// DefaultDecoderConfig returns the default decoder config which should decode
// everything that is typically used. It designed for this module and configuration
// parsing. You can also pass additional options to further configure it.
func DefaultDecoderConfig(opts ...DecodeOption) *mapstructure.DecoderConfig {
	dc := &mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.ComposeDecodeHookFunc(DefaultDecodeHooks()...),
		WeaklyTypedInput: true,
		TagName:          DefaultTag,
	}
	for _, f := range opts {
		f(dc)
	}
	return dc
}

// WithDecodeHook sets the decode hooks for this decoder
func WithDecodeHook(hooks ...mapstructure.DecodeHookFunc) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		if len(hooks) == 0 {
			dc.DecodeHook = nil
		} else if len(hooks) == 1 {
			dc.DecodeHook = hooks[0]
		} else {
			dc.DecodeHook = mapstructure.ComposeDecodeHookFunc(hooks...)
		}
	}
}

// If ErrorUnused is true, then it is an error for there to exist
// keys in the original map that were unused in the decoding process
// (extra keys).
func WithErrUnused(b bool) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.ErrorUnused = b
	}
}

// ZeroFields, if set to true, will zero fields before writing them.
// For example, a map will be emptied before decoded values are put in
// it. If this is false, a map will be merged.
func WithZeroFields(b bool) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.ZeroFields = b
	}
}

// If WeaklyTypedInput is true, the decoder will make the following
// "weak" conversions:
//
//   - bools to string (true = "1", false = "0")
//   - numbers to string (base 10)
//   - bools to int/uint (true = 1, false = 0)
//   - strings to int/uint (base implied by prefix)
//   - int to bool (true if value != 0)
//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
//     FALSE, false, False. Anything else is an error)
//   - empty array = empty map and vice versa
//   - negative numbers to overflowed uint values (base 10)
//   - slice of maps to a merged map
//   - single values are converted to slices if required. Each
//     element is weakly decoded. For example: "4" can become []int{4}
//     if the target type is an int slice.
//
func WithWeaklyTypedInput(b bool) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.WeaklyTypedInput = b
	}
}

// Squash will squash embedded structs.  A squash tag may also be
// added to an individual struct field using a tag.  For example:
//
//  type Parent struct {
//      Child `mapstructure:",squash"`
//  }
func WithSquash(b bool) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.Squash = b
	}
}

// WithTag sets the decode field tag
func WithTagName(tagName string) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.TagName = tagName
	}
}

// WithMatchName allows you to set a function to match the tag/field name with the map key.
func WithMatchName(f func(mapKey, fieldName string) bool) DecodeOption {
	return func(dc *mapstructure.DecoderConfig) {
		dc.MatchName = f
	}
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// ToSnakeCase converts to snake case.
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// MatchSnakeCaseConfig can be used to match snake case config into Go struct fields
// that do not have tags.
func MatchSnakeCaseConfig(mapKey, fieldName string) bool {
	return strings.EqualFold(mapKey, ToSnakeCase(fieldName))
}
