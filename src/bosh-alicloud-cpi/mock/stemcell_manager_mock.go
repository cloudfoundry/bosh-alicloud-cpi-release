/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type StemcellManagerMock struct {
	mc *TestContext
}

func NewStemcellManagerMock(mc TestContext) alicloud.StemcellManager {
	return StemcellManagerMock{&mc}
}

func (a StemcellManagerMock) FindStemcellById(id string) (*ecs.Image, error) {
	i, ok := a.mc.Stemcells[id]
	if !ok {
		return nil, nil
	} else {
		return i, nil
	}
}

func (a StemcellManagerMock) DeleteStemcell(id string) error {
	_, ok := a.mc.Stemcells[id]
	if !ok {
		return fmt.Errorf("DeleteImage image not exists %s", id)
	}
	delete(a.mc.Stemcells, id)
	return nil
}

func (a StemcellManagerMock) ImportImage(args *ecs.ImportImageRequest) (string, error) {
	id, image := a.mc.NewStemcell()

	image.ImageName = args.ImageName
	// ...

	return id, nil
}

func (a StemcellManagerMock) CopyImage(args *ecs.CopyImageRequest) (string, error) {
	id, image := a.mc.NewStemcell()

	image.ImageName = args.DestinationImageName
	// ...

	return id, nil
}

func (a StemcellManagerMock) OpenLocalFile(path string) (*os.File, error) {
	return nil, nil
}

func (a StemcellManagerMock) WaitForImageReady(id string) error {
	return nil
}
