// Package config provides a wrapper around YAML configs.
//
// The default name is config.yaml, with an optional overrides.yaml
// file, for local overrides.
package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	ConfigName         = "config.yaml"    // name of YAML config file
	OverridesName      = "overrides.yaml" // name of optional YAML config holding local overrides
	BasePath           = "."              // where to start looking for configs; relative to importing code
	MaxSteps      uint = 5                // maximum number of directories to step up while looking for configs

)

// MustLoad parses the YAML-encoded config and stores the result in
// the value pointed to by v.
//
// MustLoad panics if the config can't be loaded.
func MustLoad(v interface{}) {
	err := Load(v)
	if err != nil {
		panic(fmt.Errorf("FATAL: %v\n", err))
	}
}

// Load parses the YAML-encoded config and stores the result in the
// value pointed to by v.
func Load(v interface{}) error {
	err := tryLoad(ConfigName, v)
	if err != nil {
		return err
	}
	// Note: Since it's not required to have an overrides.yaml, we
	// treat a failed load as a non-error. It would be nice to log an
	// INFO message at this point to alert the caller that overrides
	// file is missing (to make the feature more discoverable), but we
	// can't use glog in case we're called from AppEngine.
	_ = tryLoad(OverridesName, v)
	return nil
}

// tryLoad parses the YAML-encoded config in file name and stores the
// result in the value pointed to by v.
//
// tryLoad steps up one directory level at a time, at most MaxSteps
// number of times, until the named config file is found.
func tryLoad(name string, v interface{}) error {
	var err error
	tries := uint(0)
	path := filepath.Join(BasePath, name)
	for tries <= MaxSteps {
		err := loadPath(path, v)
		if err == nil {
			return nil
		}
		path = filepath.Join(BasePath, strings.Repeat("../", int(tries+1)), name)
		tries += 1
	}
	return fmt.Errorf("failed to find a valid %q: %v", name, err)
}

// loadPath parses the YAML-encoded config at path and stores the
// result in the value pointed to by v.
func loadPath(path string, v interface{}) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("couldn't read config: %v", err)
	}

	err = yaml.Unmarshal(b, v)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal config: %v", err)
	}
	return nil
}
