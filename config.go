package sergeant

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/albatross-org/go-albatross/albatross"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v3"
)

// DefaultSetAll is the config for the default set, All, which contains all
var DefaultSetAll = ConfigSet{
	Name:        "All",
	Description: "This set contains all cards added to the program.",
	Color:       "#d7816a",
	Background:  "linear-gradient(315deg, #bd4f6c 0%, #d7816a 74%)",
}

// Config represents the top-level configuration for the program.
type Config struct {
	Names ConfigNames
	Sets  map[string]ConfigSet
	Store *albatross.Config
}

// ConfigSet represents the definition of a set, as specified in the config file.
type ConfigSet struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	PathsOr []string `yaml:"paths"`
	TagsOr  []string `yaml:"tags"`

	PathsAnd []string `yaml:"paths-and"`
	TagsAnd  []string `yaml:"tags-and"`

	BeforeDuration time.Duration
	AfterDuration  time.Duration
	BeforeDate     time.Time
	AfterDate      time.Time

	Color      string
	Background string
}

// AsFilter returns a ConfigSet as a filter that allows cards only if they're supposed to be in that set.
func (set ConfigSet) AsFilter() Filter {
	filters := []Filter{}

	if len(set.PathsOr) > 0 {
		filters = append(filters, FilterPaths(set.PathsOr...))
	}

	if len(set.TagsOr) > 0 {
		filters = append(filters, FilterTags(set.TagsOr...))
	}

	for _, pathAnd := range set.PathsAnd {
		filters = append(filters, FilterPaths(pathAnd))
	}

	for _, tagsAnd := range set.TagsAnd {
		filters = append(filters, FilterPaths(tagsAnd))
	}

	if set.BeforeDate != (time.Time{}) {
		filters = append(filters, FilterBeforeDate(set.BeforeDate))
	}

	if set.AfterDate != (time.Time{}) {
		filters = append(filters, FilterAfterDate(set.AfterDate))
	}

	if set.BeforeDuration != time.Duration(0) {
		filters = append(filters, FilterBeforeDuration(set.BeforeDuration))
	}

	if set.AfterDuration != time.Duration(0) {
		filters = append(filters, FilterAfterDuration(set.AfterDuration))
	}

	return FilterAND(filters...)
}

// ConfigNames lets you give friendlier names to paths to specific questions.
type ConfigNames map[string]string

// rawConfigDef is the config definition before additional processing is done on it.
// This is needed to allow the program to parse the fields instead of using YAML's default unmarshaler.
type rawConfigDef struct {
	Sets map[string]rawConfigSetDef `yaml:"sets"`

	Store *albatross.Config `yaml:"store"`
}

// LoadConfig returns the Config located at the given path. If no path is specified, the default ".config/sergeant/config.yaml" is used.
func LoadConfig(path string) (Config, error) {
	if path == "" {
		path = filepath.Join(getConfigDir(), "sergeant", "config.yaml")
	}

	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't read config located at %q: %w", path, err)
	}

	rawConfig := rawConfigDef{}

	err = yaml.Unmarshal(contentBytes, &rawConfig)
	if err != nil {
		return Config{}, fmt.Errorf("couldn't unmarshal config located at %q: %w", path, err)
	}

	config := Config{
		Sets: make(map[string]ConfigSet),
	}

	for name, rawConfigSet := range rawConfig.Sets {
		configSet, err := parseRawConfigSetDef(rawConfigSet)
		if err != nil {
			return Config{}, err
		}

		config.Sets[name] = configSet
	}

	config.Sets["all"] = DefaultSetAll
	config.Store = rawConfig.Store

	setStoreDefaults(config.Store)

	return config, nil
}

// setStoreDefaults sets the default options for a config.Store.
func setStoreDefaults(c *albatross.Config) {
	if c.DateFormat == "" {
		c.DateFormat = albatross.DefaultConfig.DateFormat
	}

	if c.TagPrefix == "" {
		c.TagPrefix = albatross.DefaultConfig.TagPrefix
	}

	// TODO: need to flip this around to c.DisableGit because it's impossible to distinguish between c.UseGit being uninitialised or
	// purposefully being set to false.
	if !c.UseGit {
		c.UseGit = true
	}

	if c.Encryption == nil {
		c.Encryption = &albatross.EncryptionConfig{}
	}

	if c.Encryption.PublicKey == "" {
		c.Encryption.PublicKey = albatross.DefaultConfig.Encryption.PublicKey
	}

	if c.Encryption.PrivateKey == "" {
		c.Encryption.PrivateKey = albatross.DefaultConfig.Encryption.PrivateKey
	}
}

// rawConfigSetDef is a definition of a set before additional processing is done on it.
// This is needed to allow the program to parse the fields such as BeforeDuration.
type rawConfigSetDef struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	PathsOr  []string `yaml:"paths"`
	TagsOr   []string `yaml:"tags"`
	PathsAnd []string `yaml:"paths-and"`
	TagsAnd  []string `yaml:"tags-and"`

	BeforeDuration string `yaml:"before-duration"`
	AfterDuration  string `yaml:"after-duration"`
	BeforeDate     string `yaml:"before-date"`
	AfterDate      string `yaml:"after-date"`

	Color      string `yaml:"color"`
	Background string `yaml:"background"`
}

// parseRawConfigSetDef turns a rawConfigSetDef into a ConfigSet.
func parseRawConfigSetDef(rawConfigSet rawConfigSetDef) (ConfigSet, error) {
	set := ConfigSet{}
	var err error

	set.Name = rawConfigSet.Name
	set.PathsAnd = rawConfigSet.PathsAnd
	set.PathsOr = rawConfigSet.PathsOr
	set.TagsAnd = rawConfigSet.TagsAnd
	set.TagsOr = rawConfigSet.TagsOr

	set.Color = rawConfigSet.Color
	set.Background = rawConfigSet.Background

	if rawConfigSet.Description == "" {
		set.Description = "This is a custom set."
	} else {
		set.Description = rawConfigSet.Description
	}

	if rawConfigSet.BeforeDate != "" {
		set.BeforeDate, err = time.Parse("2006-01-02 15:04", rawConfigSet.BeforeDate)
		if err != nil {
			return ConfigSet{}, fmt.Errorf("couldn't parse before-date %q in %q set: %w", rawConfigSet.BeforeDate, rawConfigSet.Name, err)
		}
	}

	if rawConfigSet.AfterDate != "" {
		set.AfterDate, err = time.Parse("2006-01-02 15:04", rawConfigSet.AfterDate)
		if err != nil {
			return ConfigSet{}, fmt.Errorf("couldn't parse after-date %q in %q set: %w", rawConfigSet.AfterDate, rawConfigSet.Name, err)
		}
	}

	if rawConfigSet.BeforeDuration != "" {
		set.BeforeDuration, err = time.ParseDuration(rawConfigSet.BeforeDuration)
		if err != nil {
			return ConfigSet{}, fmt.Errorf("couldn't parse before-duration %q in %q set: %w", rawConfigSet.BeforeDuration, rawConfigSet.Name, err)
		}
	}

	if rawConfigSet.AfterDuration != "" {
		set.AfterDuration, err = time.ParseDuration(rawConfigSet.AfterDuration)
		if err != nil {
			return ConfigSet{}, fmt.Errorf("couldn't parse after-duration %q in %q set: %w", rawConfigSet.AfterDuration, rawConfigSet.Name, err)
		}
	}

	return set, nil
}

// getConfigDir gets the user's configuration directory.
// TODO: At the moment, this uses $XDG_CONFIG_HOME and falls back to
// $HOME/.config which isn't cross platform.
func getConfigDir() string {
	config := os.Getenv("XDG_CONFIG_HOME")
	if config != "" {
		return config
	}

	home, err := homedir.Dir()
	if err != nil {
		panic(err) // This really shouldn't happen.
	}

	return filepath.Join(home, ".config")
}
