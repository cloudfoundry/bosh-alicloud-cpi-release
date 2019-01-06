/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
)

type TagResourceType string

const (
	TagResourceImage    = TagResourceType("image")
	TagResourceInstance = TagResourceType("instance")
	TagResourceSnapshot = TagResourceType("snapshot")
	TagResourceDisk     = TagResourceType("disk")
)

type InstanceStatus string

// Constants of InstanceStatus
const (
	Creating = InstanceStatus("Creating") // For backward compatibility
	Pending  = InstanceStatus("Pending")
	Running  = InstanceStatus("Running")
	Starting = InstanceStatus("Starting")

	Stopped  = InstanceStatus("Stopped")
	Stopping = InstanceStatus("Stopping")
	Deleted  = InstanceStatus("Deleted")
)

type EipStatus string

const (
	EipStatusAssociating   = EipStatus("Associating")
	EipStatusUnassociating = EipStatus("Unassociating")
	EipStatusInUse         = EipStatus("InUse")
	EipStatusAvailable     = EipStatus("Available")
)

type DiskStatus string

const (
	DiskStatusInUse     = DiskStatus("In_use")
	DiskStatusAvailable = DiskStatus("Available")
	DiskStatusAttaching = DiskStatus("Attaching")
	DiskStatusDetaching = DiskStatus("Detaching")
	DiskStatusCreating  = DiskStatus("Creating")
	DiskStatusReIniting = DiskStatus("ReIniting")
	DiskStatusAll       = DiskStatus("All") //Default
)

type DiskCategory string

const (
	DiskCategoryAll             = DiskCategory("all") //Default
	DiskCategoryCloud           = DiskCategory("cloud")
	DiskCategoryEphemeral       = DiskCategory("ephemeral")
	DiskCategoryEphemeralSSD    = DiskCategory("ephemeral_ssd")
	DiskCategoryCloudEfficiency = DiskCategory("cloud_efficiency")
	DiskCategoryCloudSSD        = DiskCategory("cloud_ssd")
)

type SpotStrategyType string

// Constants of SpotStrategyType
const (
	NoSpot             = SpotStrategyType("NoSpot")
	SpotWithPriceLimit = SpotStrategyType("SpotWithPriceLimit")
	SpotAsPriceGo      = SpotStrategyType("SpotAsPriceGo")
)

type ImageFormatType string

const (
	RAW = ImageFormatType("RAW")
	VHD = ImageFormatType("VHD")
)

func getSdkConfig() *sdk.Config {
	return sdk.NewConfig().
		WithMaxRetryTime(5).
		WithTimeout(time.Duration(60) * time.Second).
		WithUserAgent(getUserAgent()).
		WithGoRoutinePoolSize(10).
		WithDebug(false).
		WithHttpTransport(getTransport()).
		WithScheme("HTTPS")
}

func getTransport() *http.Transport {
	handshakeTimeout, err := strconv.Atoi(os.Getenv("TLSHandshakeTimeout"))
	if err != nil {
		handshakeTimeout = 120
	}
	return &http.Transport{
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second}
}

func getUserAgent() string {
	return "Bosh-Alicloud-Cpi"
}
