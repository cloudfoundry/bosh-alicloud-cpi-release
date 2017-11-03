/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"encoding/json"
)

type SetVMMetadataMethod struct {
	CallContext
	instances alicloud.InstanceManager
}

func NewSetVMMetadataMethod(cc CallContext, instances alicloud.InstanceManager) SetVMMetadataMethod {
	return SetVMMetadataMethod{cc, instances}
}

type MetaData map[string]interface{}

func (a SetVMMetadataMethod) SetVMMetadata(vmCID apiv1.VMCID, meta apiv1.VMMeta) error {
	bytes, err := meta.MarshalJSON()
	if err != nil {
		return a.WrapErrorf(err, "meta marshal failed")
	}
	var md MetaData
	err = json.Unmarshal(bytes, &md)
	msg := string(bytes)
	_ = msg
	if err != nil {
		return a.WrapErrorf(err, "meta unmarshal failed %s", string(bytes))
	}

	var args ecs.ModifyInstanceAttributeArgs
	args.InstanceId = vmCID.AsString()

	if s, ok := md["name"]; ok {
		args.InstanceName = s.(string)
	}

	desc := ""
	if s, ok := md["director"]; ok {
		desc += "director: " + s.(string) + "\n"
	}
	if s, ok := md["deployment"]; ok {
		desc += "deployment: " + s.(string) + "\n"
	}

	args.Description = desc
	err = a.instances.ModifyInstanceAttribute(args)

	if err != nil {
		return a.WrapErrorf(err, "set_vm_metadata failed %s", args.InstanceId)
	}

	return nil
}
