package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/egeneralov/ansible-dynamic-inventory-nspawn/internal/types"
	"go.uber.org/zap"
)

func main() {
	defer func() { _ = logger.Sync() }()
	task, err := cfg.execTask()
	if err != nil {
		logger.Fatal("failed to create exec task", zap.Error(err))
	}
	logger.Info("executing", zap.String("command", task.Command), zap.Strings("args", task.Args))
	taskResult, err := task.Execute()
	if err != nil {
		logger.Fatal("failed to execute task", zap.Error(err))
	}
	var machinesList types.List
	if err := json.Unmarshal([]byte(taskResult.Stdout), &machinesList); err != nil {
		logger.Fatal("failed to parse command stdout into types.list structure", zap.Error(err))
	}
	for _, el := range machinesList {
		logger.Debug("founded machine", zap.String("name", el.Machine), zap.String("address", el.Addresses))
	}

	if cfg.List {
		groups := machinesList.ToGroups(cfg.Bastion)
		if len(groups) == 0 {
			os.Exit(0)
		}
		jr, je := json.Marshal(groups)
		if je != nil {
			logger.Fatal("failed to marshal", zap.Error(je))
		}
		//fmt.Println(string(jr))
		if written, err := io.Copy(os.Stdout, bytes.NewBuffer(jr)); err != nil {
			logger.Fatal("failed to write stdout answer", zap.Int64("len", written), zap.Error(err))
		}
	} else if cfg.Host != "" {
		logger.Fatal("not implemented")
	} else {
		logger.Fatal("no default action available")
	}
}
