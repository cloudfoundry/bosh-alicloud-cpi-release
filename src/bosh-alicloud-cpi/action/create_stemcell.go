package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"strings"
)

type CreateStemcellMethod struct {
	runner alicloud.Runner
}

type StemcellProps struct {
	Architecture string 	`json:"architecture"`
	ContainerFormat string 	`json:"container_format"`
	Disk string				`json:"disk"`
	DiskFormat string 		`json:"disk_format"`
	Hypervisor string 		`json:"hypervisor"`
	Name string 			`json:"name"`
	OsDistro string			`json:"os_distro"`
	OsType string 			`json:"os_type"`
	RootDeviceName string 	`json:"root_device_name"`
	SourceUrl string		`json:"source_url"`
//	Version string 			`json:"version"`		TODO  sometimes string, and sometimes int
	Images map[string]interface{} 	`json:"image_id"`
}

func NewCreateStemcellMethod(runner alicloud.Runner) CreateStemcellMethod {
	return CreateStemcellMethod{runner}
}

func (a CreateStemcellMethod) CreateStemcell(imagePath string, cloudProps apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	var props StemcellProps
	err := cloudProps.As(&props)

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "BadInput for CreateStemcell %s", cloudProps)
	}

	//
	// find stemcell from manifest.MF
	region := a.runner.Config.OpenApi.RegionId
	stemcellId, err := props.FindStemcellId(region)

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
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
