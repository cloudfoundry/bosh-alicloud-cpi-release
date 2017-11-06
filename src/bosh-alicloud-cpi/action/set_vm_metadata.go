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
		args.InstanceName = normalizeInstanceName(s.(string))
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


//
// InstanceName ref https://help.aliyun.com/document_detail/25503.html
func normalizeInstanceName(s string) (string) {
	r := ""

	// can only contains [a-zA-Z0-9-_\.]
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z' || (c >= '0' && c <= '9')) || c == '-' || c == '_' || c == '.' {
			r = r + string(c)
		} else {
			r = r + "."
		}
	}

	// must start with [a-z, A-Z]
	if s[0] >= '0' && s[0] <= '9' {
		r = "i_" + r
	}

	// length in [2, 128]
	if len(r) < 2 {
		return "i_" + r
	}

	if len(r) > 128 {
		return r[0:127]
	}

	return r
}