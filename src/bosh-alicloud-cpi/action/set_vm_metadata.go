/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type SetVMMetadataMethod struct{
	config alicloud.Config
}

func NewSetVMMetadataMethod() SetVMMetadataMethod {
	return SetVMMetadataMethod{}
}

func (a SetVMMetadataMethod) SetVMMetadata(_ apiv1.VMCID, meta apiv1.VMMeta) error {
	// todo can properties be set on the container

	//client := a.config.NewEcsClient()
	//
	//
	//client.ModifyInstanceAttribute()


	return nil
}
