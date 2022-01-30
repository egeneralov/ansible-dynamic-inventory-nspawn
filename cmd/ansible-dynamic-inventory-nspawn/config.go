package main

import (
	"encoding/json"
	"errors"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CommandList string `json:"command_list" yaml:"command_list" flag:"command_list" default:"machinectl list -o json"`
	List        bool   `json:"list" yaml:"list" flag:"list" default:"true"`
	Host        string `json:"host" yaml:"host" flag:"host" default:""`
}

func (c Config) execTask() (task *execute.ExecTask, _ error) {
	if c.CommandList == "" {
		return nil, errors.New("command are empty")
	}
	task = &execute.ExecTask{}
	result := strings.Split(c.CommandList, " ")
	switch len(result) {
	case 0:
		return nil, errors.New("slice are empty")
	case 1:
		task.Command = result[0]
	default:
		task.Command = result[0]
		task.Args = result[1:]
	}
	return task, nil
}

func (c Config) Json() ([]byte, error) {
	j, je := json.Marshal(cfg)
	if je != nil {
		return nil, je
	}
	return j, nil
}

func (c Config) Yaml() ([]byte, error) {
	j, je := yaml.Marshal(cfg)
	if je != nil {
		return nil, je
	}
	return j, nil
}
