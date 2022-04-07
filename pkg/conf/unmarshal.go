package conf

import (
	"fmt"

	"github.com/knadh/koanf/maps"
	"github.com/mitchellh/mapstructure"
)

// UnmarshalConf is used to configure the unmarshaler
type UnmarshalConf struct {
	// The path in the configuration to start unmarshaling from.
	// Leave empty to unmarshal the root config structure.
	Path string
	// FlatPath interprets keys with delimeters literally instead of recursively unmarshaling structs.
	FlatPaths bool
	// Decoder config is the github.com/mitchellh/mapstructure.DecoderConfig used to umarshal
	// configuration into data structures.
	DecoderConfig *mapstructure.DecoderConfig
}

// Unmarshal configuration to dest (it should be a pointer).
// See help for UnmarshalConf for more information.
func (c *Conf) Unmarshal(dest interface{}, unmarshalConfig UnmarshalConf) error {

	// If no UnmarshalConf is specified, use the default
	if unmarshalConfig.DecoderConfig == nil {
		unmarshalConfig.DecoderConfig = DefaultDecoderConfig()
		unmarshalConfig.DecoderConfig.TagName = c.Tag
	}
	unmarshalConfig.DecoderConfig.Result = dest

	// Get the source map
	source := c.Get(unmarshalConfig.Path)
	// Flatten if requested
	if unmarshalConfig.FlatPaths {
		if f, ok := source.(map[string]interface{}); ok {
			fmp, _ := maps.Flatten(f, nil, c.Delimiter)
			source = fmp
		}
	}

	// Create the decoder and deccode into the destination.
	decoder, err := mapstructure.NewDecoder(unmarshalConfig.DecoderConfig)
	if err != nil {
		return fmt.Errorf("could not create decoder: %w", err)
	}
	return decoder.Decode(source)

}

// UnmarshalWithOpts unmarshals config into a data struct using config options.
func (c *Conf) UnmarshalWithOpts(dest interface{}, opts ...UnmarshalOption) error {
	var unmarshalConfig UnmarshalConf
	unmarshalConfig.DecoderConfig = DefaultDecoderConfig()
	for _, f := range opts {
		f(&unmarshalConfig)
	}
	return c.Unmarshal(dest, unmarshalConfig)
}

// UnmarshalOption are used to configure the Unmarshal behavior
type UnmarshalOption func(c *UnmarshalConf)

// WithPath sets the unmarshal path.
func WithPath(path string) UnmarshalOption {
	return func(c *UnmarshalConf) {
		c.Path = path
	}
}

// WithTag sets the unmarshal tag.
func WithTag(tag string) UnmarshalOption {
	return func(c *UnmarshalConf) {
		c.DecoderConfig.TagName = tag
	}
}

// WithFlatPath set unmarshaling to use flat path. See UnmarshalConf.
func WithFlatPaths(b bool) UnmarshalOption {
	return func(c *UnmarshalConf) {
		c.FlatPaths = b
	}
}

// DecoderConfig is the decoder config used to decode into the struct.
func WithDecoderOpts(opts ...DecodeOption) UnmarshalOption {
	return func(c *UnmarshalConf) {
		for _, f := range opts {
			f(c.DecoderConfig)
		}
	}
}
