package action

import (
	"bosh-alicloud-cpi/alicloud"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type CreateStemcellMethod struct {
	runner alicloud.Runner
}

func NewCreateStemcellMethod(runner alicloud.Runner) CreateStemcellMethod {
	return CreateStemcellMethod{runner}
}

func (this CreateStemcellMethod) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	// stemcell, err := a.stemcellImporter.ImportFromPath(imagePath)
	stemcellId, err := this.runner.FindStemcellId()

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	return apiv1.NewStemcellCID(stemcellId), nil
}
