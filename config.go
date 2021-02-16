package sergeant

import "time"

// DefaultSetAll is the config for the default set, All, which contains all
var DefaultSetAll = ConfigSet{
	Name:        "All",
	Description: "This set contains all cards added to the program.",
}

// ConfigSet represents the definition of a set, as specified in the config file.
type ConfigSet struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Paths []string `yaml:"paths"`
	Tags  []string `yaml:"tags"`

	BeforeDuration time.Duration
	AfterDuration  time.Duration
	BeforeDate     time.Time
	AfterDate      time.Time
}

// ConfigNames lets you give friendlier names to paths to specific questions.
type ConfigNames map[string]string

// Config represents the top-level configuration for the program.
type Config struct {
	Names ConfigNames
	Sets  map[string]ConfigSet
}

// rawConfigSet is a definition of a set before additional processing is done on it.
// This is needed to allow the program to parse the fields such as BeforeDuration.
type rawConfigSet struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	Paths []string `yaml:"paths"`
	Tags  []string `yaml:"tags"`

	BeforeDuration string `yaml:"before-duration"`
	AfterDuration  string `yaml:"after-duration"`
	BeforeDate     string `yaml:"before-date"`
	AfterDate      string `yaml:"after-date"`
}
