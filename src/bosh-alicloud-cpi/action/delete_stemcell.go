package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteStemcellMethod struct {
	config alicloud.AlicloudConfig
}

func NewDeleteStemcellMethod(config alicloud.AlicloudConfig) DeleteStemcellMethod {
	return DeleteStemcellMethod{config: config}
}

func (a DeleteStemcellMethod) DeleteStemcell(cid apiv1.StemcellCID) error {
	//stemcell, found, err := a.stemcellFinder.Find(cid)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding stemcell '%s'", cid)
	//}
	//
	//if found {
	//	err := stemcell.Delete()
	//	if err != nil {
	//		return bosherr.WrapErrorf(err, "Deleting stemcell '%s'", cid)
	//	}
	//}

	return bosherr.Error("NOT IMPLEMENTED")
}
