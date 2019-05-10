//
// (C) Copyright 2019 Intel Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// GOVERNMENT LICENSE RIGHTS-OPEN SOURCE SOFTWARE
// The Government's rights to use, modify, reproduce, release, perform, display,
// or disclose this software are subject to the terms of the Apache License as
// provided in Contract No. 8F-30005.
// Any reproduction of computer software, computer software documentation, or
// portions thereof marked with this legend must also reproduce the markings.
//

package client

import (
	"github.com/daos-stack/daos/src/control/common"
	"github.com/daos-stack/daos/src/control/log"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	daosAgentDrpcSockEnv = "DAOS_AGENT_DRPC_DIR"
	defaultRuntimeDir    = "/var/run/daos_agent"
	defaultLogFile       = "/tmp/daos_agent.log"
	defaultConfigPath    = "etc/daos.yml"
	defaultSystemName    = "daos_server"
	defaultPort          = 10000
)

// External interface provides methods to support various os operations.
type External interface {
	getenv(string) string
}

type ext struct{}

// runCommand executes command in subshell (to allow redirection) and returns
// error result.
func (e *ext) runCommand(cmd string) error {
	return exec.Command("sh", "-c", cmd).Run()
}

// getEnv wraps around os.GetEnv and implements External.getEnv().
func (e *ext) getenv(key string) string {
	return os.Getenv(key)
}

// Configuration contains all known configuration variables available to the client
type Configuration struct {
	SystemName    string   `yaml:"name"`
	AccessPoints  []string `yaml:"access_points"`
	Port          int      `yaml:"port"`
	HostList      []string `yaml:"hostlist"`
	RuntimeDir    string   `yaml:"runtime_dir"`
	HostFile      string   `yaml:"host_file"`
	Cert          string   `yaml:"cert"`
	Key           string   `yaml:"key"`
	LogFile       string   `yaml:"log_file"`
	LogFileFormat string   `yaml:"log_file_format"`
	Path          string
	ext           External
}

// newDefaultConfiguration creates a new instance of configuration struct
// populated with defaults.
func newDefaultConfiguration(ext External) Configuration {
	return Configuration{
		SystemName:   defaultSystemName,
		AccessPoints: []string{"localhost"},
		Port:         defaultPort,
		HostList:     []string{"localhost:10001"},
		RuntimeDir:   defaultRuntimeDir,
		LogFile:      defaultLogFile,
		Path:         defaultConfigPath,
		ext:          ext,
	}
}

// ProcessConfigFile loads a configuration file from the path given,
// or from the default location if none is provided.  It returns a populated
// Configuration struct based upon the default values and any config file overrides.
func ProcessConfigFile(ConfigPath string) (Configuration, error) {
	config := NewConfiguration()
	if ConfigPath != "" {
		log.Debugf("Overriding default config path with: %s", ConfigPath)
		config.Path = ConfigPath
	}

	if !filepath.IsAbs(config.Path) {
		newPath, err := common.GetAbsInstallPath(config.Path)
		if err != nil {
			return config, err
		}
		config.Path = newPath
	}

	err := config.LoadConfig()
	if err != nil {
		log.Debugf("\nUnable to read the configuration file from: '%s'\nUsing default configuration.", config.Path)
		return config, nil
	}
	log.Debugf("DAOS Client config read from %s", config.Path)

	return config, nil
}

// LoadConfig reads the configuration file specified by Configuration.Path
// and parses it.  Parsed values override any default values.
func (c *Configuration) LoadConfig() error {
	bytes, err := ioutil.ReadFile(c.Path)
	if err != nil {
		return err
	}

	if err = c.parse(bytes); err != nil {
		return err
	}

	return nil
}

// ValidateEnv reads the value of a given environment variable and compares it with
// the supplied value.  Issues a non-fatal debug message if there is a mismatch
func (c *Configuration) ValidateEnv(envVar string, value string) {
	temp := c.ext.getenv(envVar)
	if temp != value {
		log.Debugf("Environment variable '%s' has value '%s' which does not match '%s'", envVar, temp, value)
	}
}

// decodes YAML representation of Configuration struct
func (c *Configuration) parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

// NewConfiguration creates a new instance of the Configuration struct
// populated with defaults and default external interface.
func NewConfiguration() Configuration {
	return newDefaultConfiguration(&ext{})
}
