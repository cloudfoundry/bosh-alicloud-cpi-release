/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"strings"
)

type StemcellProps struct {
	Architecture string 	`json:"architecture"`
	ContainerFormat string 	`json:"container_format"`
	//Disk string				`json:"disk"`
	//DiskFormat string 		`json:"disk_format"`
	Hypervisor string 		`json:"hypervisor"`
	//Name string 			`json:"name"`
	//OsDistro string			`json:"os_distro"`
	//OsType string 			`json:"os_type"`
	//RootDeviceName string 	`json:"root_device_name"`
	//SourceUrl string		`json:"source_url"`
//	Version string 			`json:"version"`		TODO  sometimes string, and sometimes int
	Images map[string]interface{} 	`json:"image_id"`
}

type CreateStemcellMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
}

func NewCreateStemcellMethod(cc CallContext, stemcells alicloud.StemcellManager) CreateStemcellMethod {
	return CreateStemcellMethod{cc, stemcells}
}

func (a CreateStemcellMethod) CreateStemcell(imagePath string, cloudProps apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	var props StemcellProps
	err := cloudProps.As(&props)

	if err != nil {
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "BadInput for CreateStemcell %s", cloudProps)
	}

	//
	// find stemcell from manifest.MF
	region := a.Config.OpenApi.RegionId
	stemcellId, err := props.FindStemcellId(region)

	if err != nil {
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	return apiv1.NewStemcellCID(stemcellId), nil
}

func (a StemcellProps) FindStemcellId(region string) (string, error) {
	for k, v := range a.Images {
		if strings.Compare(region, string(k)) == 0 {
			s, ok := v.(string)
			if !ok {
				return "", bosherr.Errorf("Type Failed!")
			}
			return s, nil
		}
	}
	return "", bosherr.Errorf("Can't find stemcell for Region: %s", region)
}
