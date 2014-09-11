package config

import (
	"os"
	"testing"

	. "launchpad.net/gocheck"
)

func TestConfig(t *testing.T) { TestingT(t) }

type ConfigSuite struct{}

var _ = Suite(&ConfigSuite{})

type Config struct {
	Key     string
	List    []string `config:"optional"`
	Env     string
	Bool    bool
	Builtin struct {
		AnotherKey     string `config:"optional"`
		Mandatory      string
		OneMoreBuiltin struct {
			Blah string
		}
	}
}

func (s *ConfigSuite) TestConfigOk(c *C) {
	config := Config{}

	// set environment variable to test templating
	os.Setenv("SOME_ENV_KEY", "env_value")

	err := LoadConfig("configs/correct.yaml", &config)
	c.Assert(err, IsNil)
	c.Assert(config.Key, Equals, "val")
	c.Assert(config.List, DeepEquals, []string{"val1", "val2"})
	c.Assert(config.Env, Equals, "env_value")
	c.Assert(config.Builtin.AnotherKey, Equals, "anothervalue")
	c.Assert(config.Builtin.Mandatory, Equals, "onemorevalue")
	c.Assert(config.Builtin.OneMoreBuiltin.Blah, Equals, "blah")
}

func (s *ConfigSuite) TestConfigMissingRequired(c *C) {
	config := Config{}

	err := LoadConfig("configs/missing1.yaml", &config)
	c.Assert(err.Error(), Equals, "Missing required config field: Key")
}

func (s *ConfigSuite) TestConfigMissingOptional(c *C) {
	config := Config{}

	err := LoadConfig("configs/missing3.yaml", &config)
	c.Assert(err, IsNil)
}

func (s *ConfigSuite) TestConfigMissingRequiredInBuiltin(c *C) {
	config := Config{}

	err := LoadConfig("configs/missing2.yaml", &config)
	c.Assert(err.Error(), Equals, "Missing required config field: Mandatory")
}
