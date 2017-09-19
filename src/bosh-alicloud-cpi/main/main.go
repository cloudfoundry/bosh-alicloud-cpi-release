package main

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"os"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	"bosh-alicloud-cpi/action"
	"bosh-alicloud-cpi/alicloud"
)

func main() {
	logger, fs := basicDeps()
	defer logger.HandlePanic("Main")

	if len(os.Args) != 1 {
		logger.Error("main", "Usage cpi configFile")
		os.Exit(1)
	}

	configFile := os.Args[1]
	config, err := alicloud.NewConfigFromFile(configFile, fs)

	if err != nil {
		logger.Error("main", "readConfigFailed")
		os.Exit(1)
	}

	logger.Info("LoadConfig", "load Configuration from", config)

	runner := alicloud.NewRunner(logger, config)
	cpiFactory := action.NewFactory(runner)


	cli := rpc.NewFactory(logger).NewCLIWithInOut(os.Stdin, os.Stdout, cpiFactory)

	err = cli.ServeOnce()

	if err != nil {
		logger.Error("main", "Serving once %s", err)
	}

	os.Exit(1)
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	fs := boshsys.NewOsFileSystem(logger)
	//cmdRunner := boshsys.NewExecCmdRunner(logger)
	//uuidGen := boshuuid.NewGenerator()
	return logger, fs
}
