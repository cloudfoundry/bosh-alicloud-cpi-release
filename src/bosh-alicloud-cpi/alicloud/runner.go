package alicloud

import (
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cloudfoundry/bosh-utils/logger"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"strings"
	"time"
	"bosh-alicloud-cpi/registry"
)

const (
	USE_FORCE_STOP        = true
	DEFAULT_TIMEOUT       = 1200000
	DEFAULT_WAIT_INTERVAL = 2000
)

type Runner struct {
	//DiskManager
	//InstanceManager
	//NetworkManager
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
		return ecs.Deleted, bosherr.Error("Missing Instance: id=" + instid)
	}

	return inst.Status, nil
}

func (a Runner) WaitForInstanceStatus(instid string, to_status ecs.InstanceStatus) error {
	timeout := DEFAULT_TIMEOUT
	for {
		status, err := a.GetInstanceStatus(instid)

		if err != nil {
			if status == ecs.Deleted && to_status == ecs.Deleted {
				return nil
			}

			return err
		}

		if status == to_status {
			return nil
		}

		if timeout > 0 {
			timeout -= DEFAULT_WAIT_INTERVAL
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

func (a Runner) GetHttpRegistryClient() (registry.Client) {
	r := a.Config.Registry

	if strings.Compare("", r.Host) == 0 {
		//
		// first start need skip this operation
		return nil
	}

	clientOptions := registry.ClientOptions {
		Protocol: r.Protocol,
		Host: r.Host,
		Port: r.Port,
		Username: r.User,
		Password: r.Password,
	}

	client := registry.NewHTTPClient(clientOptions, a.Logger)
	return client
}

func (a Runner) GetDisks(instid string) ([]ecs.DiskItemType, error) {
	client := a.NewClient()
	args := ecs.DescribeDisksArgs {
		RegionId: common.Region(a.Config.OpenApi.RegionId),
		InstanceId: instid,
	}
	disks, _, err := client.DescribeDisks(&args)
	return disks, err
}

func (a Runner) AttachDisk(instid string, diskid string) (error) {
	client := a.NewClient()
	var args ecs.AttachDiskArgs
	args.InstanceId = instid
	args.DiskId = diskid
	err := client.AttachDisk(&args)
	return err
}

func (a Runner) GetDiskStatus(diskid string) (*ecs.DiskItemType, error) {
	client := a.NewClient()
	args := ecs.DescribeDisksArgs {
		RegionId: common.Region(a.Config.OpenApi.RegionId),
		DiskIds: []string { diskid, },
	}
	disks, _, err := client.DescribeDisks(&args)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "GetDisk() Failed %s", args)
	}
	if len(disks) == 0 {
		return nil, nil
	}
	return &disks[0], nil
}

func (a Runner) WaitForDiskStatus(diskid string, toStatus ecs.DiskStatus) (string, error) {
	timeout := DEFAULT_TIMEOUT
	for {
		disk, err := a.GetDiskStatus(diskid)

		if err != nil {
			return "", err
		}

		if disk.Status == toStatus {
			//
			// 如果非普通云盘，需要去除x字母，如: xvdb -> vdb
			// if not normal Cloud need trim first x: xvdc -> vcd
			device := disk.Device
			if device[5] == 'x' {
				device = "/dev/" + string(device[6:])
			}

			return device, nil
		}

		if timeout > 0 {
			timeout -= 1000
			time.Sleep(time.Duration(DEFAULT_WAIT_INTERVAL) * time.Millisecond)
		} else {
			return "", bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
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
