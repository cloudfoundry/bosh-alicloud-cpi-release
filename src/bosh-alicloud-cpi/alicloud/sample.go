package alicloud

import (
bosherr "github.com/cloudfoundry/bosh-utils/errors"
"cppforlife/bosh-warden-cpi-release/src/github.com/cppforlife/bosh-cpi-go/apiv1"
"github.com/denverdino/aliyungo/ecs"
"fmt"
)

type CreateStemcellMethod struct {
	//stemcellImporter bwcstem.Importer
}

const ACCESS_KEY_ID = "***REMOVED***"
const ACCESS_KEY_SECRET = "***REMOVED***"



func (a CreateStemcellMethod) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	// stemcell, err := a.stemcellImporter.ImportFromPath(imagePath)

	client := ecs.NewClient(ACCESS_KEY_ID, ACCESS_KEY_SECRET)

	regions, err := client.DescribeRegions()

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	fmt.Print(regions)

	return apiv1.NewStemcellCID("foo-id"), nil
}