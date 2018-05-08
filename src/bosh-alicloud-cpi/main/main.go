/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package main

import (
	"bosh-alicloud-cpi/action"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var configFile = flag.String("configFile", "", `cpi -configFile=/path/to/configuration_file.json`)

func main() {
	logger, fs := basicDeps()
	defer logger.HandlePanic("Main")

	flag.Parse()
	config, err := alicloud.NewConfigFromFile(*configFile, fs)

	if err != nil {
		logger.Error("main", "read config failed %s", err)
		os.Exit(1)
	}

	logger.Info("CONFIG", "load Configuration from %s: %s", configFile, config)

	caller := action.NewCaller(config, logger)

	input, err := ioutil.ReadAll(os.Stdin)
	resp := caller.Run(input)
	output, _ := json.Marshal(resp)

	os.Stdout.Write(output)

	if err != nil {
		logger.Error("main", "Serving once %s", err)
	}

	os.Exit(0)
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	fs := boshsys.NewOsFileSystem(logger)
	//cmdRunner := boshsys.NewExecCmdRunner(logger)
	//uuidGen := boshuuid.NewGenerator()
	return logger, fs
}
