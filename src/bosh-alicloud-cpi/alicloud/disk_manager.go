package alicloud

import "github.com/denverdino/aliyungo/ecs"

type DiskManager interface {
	GetDisks(instCid string) ([]ecs.DiskItemType, error)
	CreateDisk(args ecs.CreateDiskArgs) (string, error)
	AttachDisk(instCid string, diskCid string) (error)
	DetachDisk(instCid string, diskCid string) (error)
	DeleteDisk(diskCid string) (error)
	GetDiskStatus(diskCid string) (ecs.DiskStatus, error)
	WaitDiskStatus(diskCid string, status ecs.DiskStatus) (error)
}