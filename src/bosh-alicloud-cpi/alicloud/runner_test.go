package alicloud

import (
	"testing"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"os"
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/action"
)

var testInstId = "i-2zefl7hfr7yb97ni5skw"

func TestGetVmStatus(t *testing.T) {
	config, _ := alicloud.NewConfigFromBytes(action.TestConfig)
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)

	status, err := runner.GetInstanceStatus(testInstId)

	if err != nil {
		t.Error("VmStatus", "TestGetVmStatus failed, ", err)
	}

	logger.Info("VmStatus:", "Status %s", status)
	t.Log(status)
}

func TestStopStartVm(t *testing.T) {
	config, _ := alicloud.NewConfigFromBytes(action.TestConfig)
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)

	t.Log("Stopping VM ...")
	err := runner.StopInstance(testInstId)
	if err != nil {
		t.Error("VmStatus", "Stoping failed, ", err)
	}
	err = runner.WaitForInstanceStatus(testInstId, ecs.Stopped)
	if err != nil {
		t.Error("VmStatus", "WaitForVm failed, ", err)
	}

	t.Log("Stopping VM Done")


	t.Log("Stopping VM ...")
	err = runner.StartInstance(testInstId)
	if err != nil {
		t.Error("VmStatus", "Stoping failed, ", err)
	}
	err = runner.WaitForInstanceStatus(testInstId, ecs.Stopped)
	if err != nil {
		t.Error("VmStatus", "WaitForVm failed, ", err)
	}

	t.Log("Stopping VM Done")
}
