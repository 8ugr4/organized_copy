package pkg

import (
	"go.yaml.in/yaml/v4"
	"os"
)

type Rule struct {
	Category     string   `yaml:"category"`
	Separate     []string `yaml:"separate"`
	Extensions   []string `yaml:"extensions,omitempty"`
	NameContains []string `yaml:"name_contains,omitempty"`
	Sort         string   `yaml:"sort,omitempty"` // month&year is only possible options
}

type Override struct {
	Priority []string `yaml:"priority_order,omitempty"`
}

type Config struct {
	Rules    []Rule   `yaml:"rules"`
	Override Override `yaml:"override"`
}

func ReadCategories(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (r Rule) SeparateExists() bool {
	return len(r.Separate) > 0
}
