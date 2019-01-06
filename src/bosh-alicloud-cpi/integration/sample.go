/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package integration

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type CreateStemcellMethod struct {
	//stemcellImporter bwcstem.Importer
}

const ACCESS_KEY_ID = "***your key***"
const ACCESS_KEY_SECRET = "***you secret***"
const REGION_ID = "cn-hangzhou"

func (a CreateStemcellMethod) CreateStemcell(imagePath string, _ apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	// stemcell, err := a.stemcellImporter.ImportFromPath(imagePath)

	client, err := ecs.NewClientWithOptions(REGION_ID, getSdkConfig(), credentials.NewAccessKeyCredential(ACCESS_KEY_ID, ACCESS_KEY_SECRET))
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Initiating ECS Client in '%s' got an error.", REGION_ID)
	}

	args := ecs.CreateDescribeRegionsRequest()
	regions, err := client.DescribeRegions(args)

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	fmt.Print(regions)

	return apiv1.NewStemcellCID("foo-id"), nil
}

func getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(5).
		WithUserAgent("Bosh-Alicloud-Cpi").
		WithGoRoutinePoolSize(10).
		WithDebug(false)
}
