/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type Factory struct {
	cc CallContext
	services Services
}

type CPI struct {
	InfoMethod

	CreateStemcellMethod
	DeleteStemcellMethod

	CreateVMMethod
	DeleteVMMethod
	HasVMMethod
	RebootVMMethod
	SetVMMetadataMethod
	GetDisksMethod

	CreateDiskMethod
	DeleteDiskMethod
	AttachDiskMethod
	DetachDiskMethod
	HasDiskMethod
}

func NewFactory(cc CallContext, services Services) (Factory) {
	return Factory{
		cc: cc,
		services: services,
	}
}

func (f Factory) New(_ apiv1.CallContext) (apiv1.CPI, error) {
	cc := f.cc
	ss := f.services
	cpi := CPI {
		NewInfoMethod(),

		NewCreateStemcellMethod(cc, ss.Stemcells),
		NewDeleteStemcellMethod(cc, ss.Stemcells),

		NewCreateVMMethod(cc, ss.Stemcells, ss.Instances, ss.Disks, ss.Networks, ss.Registry),
		NewDeleteVMMethod(cc, ss.Instances),
		NewHasVMMethod(cc, ss.Instances),
		NewRebootVMMethod(cc, ss.Instances),
		NewSetVMMetadataMethod(),

		NewGetDisksMethod(cc, ss.Disks),
		NewCreateDiskMethod(cc, ss.Disks),
		NewDeleteDiskMethod(cc, ss.Disks),
		NewAttachDiskMethod(cc, ss.Disks, ss.Registry),
		NewDetachDiskMethod(cc, ss.Disks, ss.Registry),
		NewHasDiskMethod(cc, ss.Disks),
	}

	return cpi, nil
}
