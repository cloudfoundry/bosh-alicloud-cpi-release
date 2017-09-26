package alicloud

import (
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"strings"
	"time"
)

const (
	USE_FORCE_STOP        = true
	DEFAULT_TIMEOUT       = 1200000
	DEFAULT_WAIT_INTERVAL = 500
)

type Runner struct {
	Logger logger.Logger
	Config Config
}

func NewRunner(logger logger.Logger, config Config) Runner {
	return Runner{
		Logger: logger,
		Config: config,
	}
}

func (a Runner) NewClient() *ecs.Client {
	return ecs.NewClient(a.Config.OpenApi.AccessKeyId, a.Config.OpenApi.AccessKeySecret)
}

func (a Runner) FindStemcellId() (string, error) {
	c := a.Config
	for _, region := range c.OpenApi.Regions {
		if strings.Compare(region.Name, c.OpenApi.RegionId) == 0 {
			return region.ImageId, nil
		}
	}
	return "", fmt.Errorf("Unknown Region")
}

func (a Runner) GetInstance(instid string) (*ecs.InstanceAttributesType, error) {
	client := a.NewClient()

	var args ecs.DescribeInstancesArgs
	args.RegionId = common.Region(a.Config.OpenApi.RegionId)
	args.InstanceIds = "[\"" + instid + "\"]"

	insts, _, err := client.DescribeInstances(&args)

	if err != nil {
		return nil, err
	}

	if len(insts) == 0 {
		return nil, nil
	}

	return &insts[0], nil
}

func (a Runner) GetInstanceStatus(instid string) (ecs.InstanceStatus, error) {
	inst, err := a.GetInstance(instid)

	if err != nil {
		return "", err
	}

	if inst == nil {
		return "", bosherr.Error("Missing Instance: id=" + instid)
	}

	return inst.Status, nil
}

func (a Runner) WaitForInstanceStatus(instid string, to_status ecs.InstanceStatus) error {
	timeout := DEFAULT_TIMEOUT
	for {
		status, err := a.GetInstanceStatus(instid)

		if err != nil {
			return err
		}

		if status == to_status {
			return nil
		}

		if timeout > 0 {
			timeout -= 1000
			time.Sleep(time.Duration(DEFAULT_WAIT_INTERVAL) * time.Millisecond)
		} else {
			return bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
}

func (a Runner) StopInstance(instid string) error {
	client := a.NewClient()
	return client.StopInstance(instid, USE_FORCE_STOP)
}

func (a Runner) StartInstance(instid string) error {
	client := a.NewClient()
	return client.StartInstance(instid)
}

func (a Runner) RebootInstance(instid string) error {
	client := a.NewClient()
	return client.RebootInstance(instid, USE_FORCE_STOP)
}

//func (a Runner) HasDisk(diskid string) (error) {
//	client := a.NewClient()
//	return client.DescribeDisks()
//}

//func (a Runner) CreateDisk() {
//	client := a.NewClient()
//	var args ecs.DescribeVSwitchesArgs;
//	args.RegionId
//	client.DescribeVSwitches()
//
//}
