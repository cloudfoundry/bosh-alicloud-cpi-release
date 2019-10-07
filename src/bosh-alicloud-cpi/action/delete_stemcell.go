/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
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
	if err := a.stemcells.DeleteStemcell(cid.AsString()); err != nil {
		return bosherr.WrapErrorf(err, "Deleting stemcell '%s'", cid)
	}

	return nil
}
