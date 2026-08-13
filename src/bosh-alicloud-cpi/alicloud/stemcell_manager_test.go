/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// fakeEcsImageClient is a hand-rolled EcsImageClient that records calls and lets
// each test drive DescribeImages/DeleteImage behavior. It records the order of
// image-vs-snapshot deletions so ordering can be asserted.
type fakeEcsImageClient struct {
	// describeImages is invoked for every DescribeImages call (initial lookup
	// plus each WaitForImageDeleted poll). callNum starts at 1.
	describeImages func(callNum int) (*ecs.DescribeImagesResponse, error)

	deleteImageErr error

	callLog          []string
	describeCalls    int
	deletedSnapshots []string
}

func (f *fakeEcsImageClient) DescribeImages(request *ecs.DescribeImagesRequest) (*ecs.DescribeImagesResponse, error) {
	f.describeCalls++
	return f.describeImages(f.describeCalls)
}

func (f *fakeEcsImageClient) DeleteImage(request *ecs.DeleteImageRequest) (*ecs.DeleteImageResponse, error) {
	f.callLog = append(f.callLog, "DeleteImage:"+request.ImageId)
	if f.deleteImageErr != nil {
		return nil, f.deleteImageErr
	}
	return &ecs.DeleteImageResponse{}, nil
}

func (f *fakeEcsImageClient) DeleteSnapshot(request *ecs.DeleteSnapshotRequest) (*ecs.DeleteSnapshotResponse, error) {
	f.callLog = append(f.callLog, "DeleteSnapshot:"+request.SnapshotId)
	f.deletedSnapshots = append(f.deletedSnapshots, request.SnapshotId)
	return &ecs.DeleteSnapshotResponse{}, nil
}

func (f *fakeEcsImageClient) ImportImage(request *ecs.ImportImageRequest) (*ecs.ImportImageResponse, error) {
	return nil, nil
}

func (f *fakeEcsImageClient) CopyImage(request *ecs.CopyImageRequest) (*ecs.CopyImageResponse, error) {
	return nil, nil
}

func (f *fakeEcsImageClient) ModifyImageAttribute(request *ecs.ModifyImageAttributeRequest) (*ecs.ModifyImageAttributeResponse, error) {
	return nil, nil
}

// imageResponse builds a DescribeImagesResponse for an image with the given
// backing snapshot IDs.
func imageResponse(imageId string, snapshotIds ...string) *ecs.DescribeImagesResponse {
	img := ecs.Image{ImageId: imageId}
	for _, ssid := range snapshotIds {
		img.DiskDeviceMappings.DiskDeviceMapping = append(
			img.DiskDeviceMappings.DiskDeviceMapping,
			ecs.DiskDeviceMapping{SnapshotId: ssid},
		)
	}
	resp := &ecs.DescribeImagesResponse{}
	resp.Images.Image = []ecs.Image{img}
	return resp
}

var _ = Describe("StemcellManagerImpl.DeleteStemcell", func() {
	var (
		fake    *fakeEcsImageClient
		manager StemcellManagerImpl
	)

	BeforeEach(func() {
		fake = &fakeEcsImageClient{}
		manager = StemcellManagerImpl{
			logger: boshlog.NewLogger(boshlog.LevelNone),
			newClient: func(region string) (EcsImageClient, error) {
				return fake, nil
			},
		}
	})

	It("deletes backing snapshots after the image is confirmed gone", func() {
		// Call 1 (initial lookup) returns the image with two snapshots.
		// Subsequent calls (WaitForImageDeleted poll) report it gone.
		fake.describeImages = func(callNum int) (*ecs.DescribeImagesResponse, error) {
			if callNum == 1 {
				return imageResponse("m-image", "s-snap1", "s-snap2"), nil
			}
			return &ecs.DescribeImagesResponse{}, nil
		}

		err := manager.DeleteStemcell("m-image")
		Expect(err).NotTo(HaveOccurred())

		Expect(fake.deletedSnapshots).To(ConsistOf("s-snap1", "s-snap2"))
		// The image must be deleted before any snapshot.
		Expect(fake.callLog[0]).To(Equal("DeleteImage:m-image"))
		Expect(fake.callLog).To(ContainElements("DeleteSnapshot:s-snap1", "DeleteSnapshot:s-snap2"))
	})

	It("does not delete snapshots when DeleteImage fails", func() {
		fake.describeImages = func(callNum int) (*ecs.DescribeImagesResponse, error) {
			return imageResponse("m-image", "s-snap1"), nil
		}
		fake.deleteImageErr = fmt.Errorf("DeleteImage: forbidden")

		err := manager.DeleteStemcell("m-image")
		Expect(err).To(HaveOccurred())
		Expect(fake.deletedSnapshots).To(BeEmpty())
	})

	It("skips snapshot deletion when the image is not confirmed deleted", func() {
		// Initial lookup succeeds; after DeleteImage the poll keeps erroring
		// (image still present / API error), so the wait fails.
		fake.describeImages = func(callNum int) (*ecs.DescribeImagesResponse, error) {
			if callNum == 1 {
				return imageResponse("m-image", "s-snap1"), nil
			}
			return nil, fmt.Errorf("DescribeImages: throttled")
		}

		err := manager.DeleteStemcell("m-image")
		Expect(err).NotTo(HaveOccurred())
		// No snapshot deletion attempted against a still-live image.
		Expect(fake.deletedSnapshots).To(BeEmpty())
	})

	It("is idempotent when the image does not exist", func() {
		fake.describeImages = func(callNum int) (*ecs.DescribeImagesResponse, error) {
			return &ecs.DescribeImagesResponse{}, nil
		}

		err := manager.DeleteStemcell("m-missing")
		Expect(err).NotTo(HaveOccurred())
		Expect(fake.callLog).To(BeEmpty())
	})
})
