package conf

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/providers/structs"
)

const DefaultTag = "conf"
const DefaultDelimiter = "."

// Conf is a wrapper around a koanf.Koanf to give it load and unmarshal functionality.
type Conf struct {
	*koanf.Koanf
	Opts
}

// Opts allows overriding the default tag and delimiters.
type Opts struct {
	Delimiter string
	Tag       string
}

// C is the optional global configuration using "." as the delimeter.
var C = New()

// New returns a new Conf with the default options.
func New() *Conf {
	return NewWithOpts(Opts{
		Delimiter: DefaultDelimiter,
		Tag:       DefaultTag,
	})

}

// NewWithOpts returns a new Conf instance with custom options.
func NewWithOpts(opts Opts) *Conf {
	return &Conf{
		Koanf: koanf.New(opts.Delimiter),
		Opts:  opts,
	}
}

// ParserFunc is a option function for loading different types of config.
type ParserFunc func(*Conf) error

// Parse is the all purpose wrapper to parse configuration from a multitude of places.
// Config sources are provided via ParserFuncs.
func (c *Conf) Parse(parsers ...ParserFunc) error {
	for _, parser := range parsers {
		if err := parser(c); err != nil {
			return err
		}
	}
	return nil
}

// WithMap is a ParserFunc to leverage a map to load configuration.
// See ParseMap for more information.
func WithMap(config map[string]interface{}) ParserFunc {
	return func(c *Conf) error {
		return c.ParseMap(config)
	}
}

// WithFile parses a configuration from a file.
// See ParseFile for information on file types supported.
func WithFile(configFile string) ParserFunc {
	return func(c *Conf) error {
		return c.ParseFile(configFile)
	}
}

// WithBytes parses a configuration from byte slice.
// See ParseBytes for information on formats supported.
func WithBytes(b []byte, format string) ParserFunc {
	return func(c *Conf) error {
		return c.ParseBytes(b, format)
	}
}

// WithStruct parses a configuration from a struct.
// See ParseStruct for information on arguments.
func WithStruct(in interface{}) ParserFunc {
	return func(c *Conf) error {
		return c.ParseStruct(in)
	}
}

// WithStructWithTag parses a configuration from a struct.
// See ParseStructWithTag for information on arguments.
func WithStructWithTag(in interface{}, tag string) ParserFunc {
	return func(c *Conf) error {
		return c.ParseStructWithTag(in, tag)
	}
}

// WithEnv parses configuration from environment variable.
// See ParseEnv for more information. Environment variables
// can only override existing configuration values.
func WithEnv() ParserFunc {
	return func(c *Conf) error {
		return c.ParseEnv()
	}
}

// ParseMap parses configuration from a map[string]interface{} and is handy for
// passing in defaults. If config is nil, it is ignored.
func (c *Conf) ParseMap(config map[string]interface{}) error {
	// If config is nil, just ignore it.
	if config == nil {
		return nil
	}
	return c.Load(confmap.Provider(config, "."), nil)
}

// ParseFile loads configuration from a file. It supports yaml, json and toml.
// The type is inferred from configFile extension. If configFile is an empty
// string the file is ignored.
func (c *Conf) ParseFile(configFile string) error {
	// If configFile is empty, just skip it.
	if configFile == "" {
		return nil
	}
	return c.ParseProvider(file.Provider(configFile), filepath.Ext(configFile))
}

// ParseBytes loads configuration from bytes. It supports yaml, json and toml.
// The format must be supplied. If buf is empty is it ignored.
func (c *Conf) ParseBytes(b []byte, format string) error {
	// If empty, just skip it.
	if len(b) == 0 {
		return nil
	}
	return c.ParseProvider(rawbytes.Provider(b), format)
}

// ParseProvider is a helper that takes a koanf provider and format and
// parses configuration from it..
func (c *Conf) ParseProvider(p koanf.Provider, format string) error {

	switch format {
	case "yaml", ".yaml", ".yml":
		return c.Load(p, yaml.Parser())
	case "json", ".json":
		return c.Load(p, json.Parser())
	case "toml", ".toml":
		return c.Load(p, toml.Parser())
	}
	return fmt.Errorf("unknown config format %s", format)
}

// ParseStruct loads configuration from a struct. If it's nil, it's ignored.
// It will follow any tags configured on the struct.
func (c *Conf) ParseStruct(in interface{}) error {
	return c.ParseStructWithTag(in, c.Tag)
}

// ParseStruct loads configuration from a struct. If it's nil, it's ignored.
// This allows you to set what tag to look at on the struct.
func (c *Conf) ParseStructWithTag(in interface{}, tag string) error {
	if in == nil {
		return nil
	}
	return c.Load(structs.Provider(in, tag), nil)
}

// WithEnv parses configuration values from environment variables. It is only
// useful for overriding values that are already present in the configuration.
// By default it looks for an environment variable that is all caps with periods
// replaced by underscores to override. database.user_password would be overridden
// by an environment variable called DATABASE_USER_PASSWORD.
func (c *Conf) ParseEnv() error {
	// All underscores in environment variables to dots
	envReplacer := strings.NewReplacer("_", ".")
	// Build a map of existing config items with all underscores replaced with dots so `thing.that_value` can
	// be replaced by environment variable THING_THAT_VALUE instead of it trying to replace `thing.that.value`
	envLookup := make(map[string]string) //
	for _, key := range c.Keys() {
		envLookup[envReplacer.Replace(key)] = key
	}
	// Load the environment variables, compare to our lookup of existing values and set override value
	return c.Load(env.ProviderWithValue("", ".", func(key string, value string) (string, interface{}) {
		// Convert environment variable to lower case and change underscore to dot.
		key = envReplacer.Replace(strings.ToLower(key))
		if replacement, found := envLookup[key]; found {
			// Check the existing type of the variable, and allow modifying.
			switch c.Get(replacement).(type) {
			case []interface{}, []string: // If existing value is string slice, split on space.
				return replacement, strings.Split(value, " ")
			}
			return replacement, value
		}
		return "", nil // No existing variable, skip it
	}), nil)
}
