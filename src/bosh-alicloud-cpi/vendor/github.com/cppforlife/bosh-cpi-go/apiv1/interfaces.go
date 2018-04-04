package apiv1

type CPIFactory interface {
	New(CallContext) (CPI, error)
}

type CallContext interface {
	As(interface{}) error
}

type CPI interface {
	Info() (Info, error)
	Stemcells
	VMs
	Disks
	Snapshots
}

type Info struct {
	StemcellFormats []string `json:"stemcell_formats"`
}
