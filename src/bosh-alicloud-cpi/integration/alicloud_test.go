/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	. "github.com/onsi/ginkgo"
)

var _ = Describe("integration:alicloud", func() {
	It("can run basic", func() {

	})

	//It("get vm status", func() {
	//
	//	config, _ := NewConfigFromBytes(TestConfig)
	//	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	//	runner := NewRunner(logger, config)
	//
	//	status, err := runner.GetInstanceStatus(testInstId)
	//
	//	if err != nil {
	//		GinkgoT().Error("VmStatus", "TestGetVmStatus failed, ", err)
	//	}
	//
	//	logger.Info("VmStatus:", "Status %s", status)
	//	GinkgoT().Log(status)
	//})
	//It("stop start vm", func() {
	//
	//	config, _ := NewConfigFromBytes(TestConfig)
	//	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	//	runner := NewRunner(logger, config)
	//	GinkgoT().Log("Stopping VM ...")
	//	err := runner.StopInstance(testInstId)
	//	if err != nil {
	//		GinkgoT().Error("VmStatus", "Stoping failed, ", err)
	//	}
	//	err = runner.WaitForInstanceStatus(testInstId, ecs.Stopped)
	//	if err != nil {
	//		GinkgoT().Error("VmStatus", "WaitForVm failed, ", err)
	//	}
	//	GinkgoT().Log("Stopping VM Done")
	//	GinkgoT().Log("Starting VM ...")
	//	err = runner.StartInstance(testInstId)
	//	if err != nil {
	//		GinkgoT().Error("VmStatus", "Starting failed, ", err)
	//	}
	//	err = runner.WaitForInstanceStatus(testInstId, ecs.Running)
	//	if err != nil {
	//		GinkgoT().Error("VmStatus", "WaitForVm failed, ", err)
	//	}
	//	GinkgoT().Log("Stopping VM Done")
	//})
})
var testInstId = "i-2zefl7hfr7yb97ni5skw"
