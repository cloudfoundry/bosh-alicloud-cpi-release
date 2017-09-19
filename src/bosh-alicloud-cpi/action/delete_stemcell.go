package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteStemcellMethod struct {
	runner alicloud.Runner
}

func NewDeleteStemcellMethod(runner alicloud.Runner) DeleteStemcellMethod {
	return DeleteStemcellMethod{runner}
}

func (a DeleteStemcellMethod) DeleteStemcell(cid apiv1.StemcellCID) error {
	//
	// NOTHING TO DELETE


	//stemcell, found, err := a.runner.FindStemcellId(cid)
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

	return nil
}
