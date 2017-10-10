package alicloud

import "github.com/denverdino/aliyungo/ecs"

type InstanceManager interface {
	CreateInstance(args ecs.CreateInstanceArgs) (string, error)
	DeleteInstance(cid string) (error)
	StartInstance(cid string) (error)
	StopInstance(cid string) (error)
	GetInstanceStatus(cid string) (ecs.InstanceStatus, error)
	WaitInstanceStatus(cid string, status ecs.InstanceStatus) (error)
}