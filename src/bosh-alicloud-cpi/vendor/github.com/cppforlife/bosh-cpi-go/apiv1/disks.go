package apiv1

type Disks interface {
	CreateDisk(int, DiskCloudProps, *VMCID) (DiskCID, error)
	DeleteDisk(DiskCID) error

	AttachDisk(VMCID, DiskCID) error
	DetachDisk(VMCID, DiskCID) error
	SetDiskMetadata(DiskCID, DiskMeta) error

	HasDisk(DiskCID) (bool, error)
	ResizeDisk(DiskCID, int) error
}

type DiskCloudProps interface {
	As(interface{}) error
	_final() // interface unimplementable from outside
}

type DiskCID struct {
	cloudID
}

type DiskMeta struct {
	cloudKVs
}

func NewDiskCID(cid string) DiskCID {
	if cid == "" {
		panic("Internal inconsistency: Disk CID must not be empty")
	}
	return DiskCID{cloudID{cid}}
}
