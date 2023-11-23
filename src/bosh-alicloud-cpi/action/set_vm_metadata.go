/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
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
		return bosherr.WrapErrorf(err, "convert meta data failed %v", meta)
	}

	instCid := vmCID.AsString()

	tags := make(map[string]string)
	name := ""
	for k, v := range md {
		if k == "name" {
			name = normalizeName(v.(string), "i_")
		}
		// There should support manifest tags feature: https://bosh.io/docs/manifest-v2/#tags
		//} else if k == "deployment" || k == "director" || k == "index" || k == "instance_group" || k == "job" {
		//	tk := normalizeTag(k)
		//	if tk != "" {
		//		tags[tk] = normalizeTag(v.(string))
		//	}
		if normalizeTag(k) != "" {
			tags[normalizeTag(k)] = normalizeTag(v.(string))
		}
	}

	err = a.instances.ModifyInstanceAttribute(instCid, name, "")
	if err != nil {
		return bosherr.WrapErrorf(err, "ModifyInstanceAttribute %s failed", instCid)
	}

	err = a.instances.AddTags(instCid, tags)
	if err != nil {
		return bosherr.WrapErrorf(err, "AddTags %v to %s failed", instCid, tags, instCid)
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

// InstanceName ref https://help.aliyun.com/document_detail/25503.html
func normalizeName(s string, prefix string) string {
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

// Tag.Key Tag.Name
// ref https://help.aliyun.com/document_detail/25616.html
func normalizeTag(s string) string {
	if strings.HasPrefix(s, "aliyun") || strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		s = "_" + s
	}

	if len(s) > 128 {
		s = string(s[0:127])
	}

	return s
}
