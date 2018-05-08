/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type DeleteStemcellMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
}

func NewDeleteStemcellMethod(cc CallContext, stemcells alicloud.StemcellManager) DeleteStemcellMethod {
	return DeleteStemcellMethod{cc, stemcells}
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
