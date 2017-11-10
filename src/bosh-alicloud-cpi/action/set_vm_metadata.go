/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"fmt"
)

type SetVMMetadataMethod struct {
	CallContext
	instances alicloud.InstanceManager
}

func NewSetVMMetadataMethod(cc CallContext, instances alicloud.InstanceManager) SetVMMetadataMethod {
	return SetVMMetadataMethod{cc, instances}
}

func (a SetVMMetadataMethod) SetVMMetadata(vmCID apiv1.VMCID, meta apiv1.VMMeta) error {
	md, err := convertMetaData(meta)
	if err != nil {
		return a.WrapErrorf(err, "convert meta data failed %v", meta)
	}

	instCid := vmCID.AsString()
	name := ""
	if s, ok := md["name"]; ok {
		name = normalizeName(s.(string), "i_")
	}

	desc := ""
	if s, ok := md["director"]; ok {
		desc += "director: " + s.(string) + "\n"
	}
	if s, ok := md["deployment"]; ok {
		desc += "deployment: " + s.(string) + "\n"
	}

	err = a.instances.ModifyInstanceAttribute(instCid, name, desc)
	if err != nil {
		return a.WrapErrorf(err, "set_vm_metadata failed %s", instCid)
	}

	return nil
}

type MetaData map[string]interface{}
type MetaInput interface {
	MarshalJSON() ([]byte, error)
}

func convertMetaData(input MetaInput) (MetaData, error) {
	var r MetaData
	bytes, err := input.MarshalJSON()
	if err != nil {
		return r, fmt.Errorf("meta marshal failed %s", err.Error())
	}
	err = json.Unmarshal(bytes, &r)
	if err != nil {
		return r, fmt.Errorf("meta unmarshal failed %s", string(bytes))
	}
	return r, nil
}
//
// InstanceName ref https://help.aliyun.com/document_detail/25503.html
func normalizeName(s string, prefix string) (string) {
	r := ""

	// can only contains [a-zA-Z0-9-_\.]
	for _, c := range s {
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' {
			r = r + string(c)
		} else {
			r = r + "."
		}
	}

	// must start with [a-z, A-Z]
	c := s[0]
	if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')) {
		r = prefix + r
	}

	// length in [2, 128]
	if len(r) < 2 {
		return prefix + r
	}

	if len(r) > 128 {
		return r[0:127]
	}

	return r
}