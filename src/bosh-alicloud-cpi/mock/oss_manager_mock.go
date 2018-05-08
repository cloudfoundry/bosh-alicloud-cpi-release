package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssManagerMock struct {
	mc *TestContext
}

func NewOssManagerMock(mc TestContext) alicloud.OssManager {
	return OssManagerMock{&mc}
}

func (a OssManagerMock) CreateBucket(name string, options ...oss.Option) error {
	_, bucket := a.mc.NewBucket(name)
	bucket.BucketName = name
	// ...

	return nil
}

func (a OssManagerMock) DeleteBucket(name string) error {
	_, ok := a.mc.Buckets[name]
	if !ok {
		return fmt.Errorf("DeleteBucket bucket not exists %s", name)
	}
	delete(a.mc.Buckets, name)
	return nil
}

func (a OssManagerMock) GetBucket(name string) (*oss.Bucket, error) {
	b, ok := a.mc.Buckets[name]
	if !ok {
		return nil, nil
	} else {
		return b, nil
	}
}

func (a OssManagerMock) UploadFile(
	bucket oss.Bucket, objectKey, filePath string, partSize int64, options ...oss.Option) error {
	a.mc.NewObject(objectKey, filePath)
	return nil
}

func (a OssManagerMock) DeleteObject(bucket oss.Bucket, name string) error {
	_, ok := a.mc.OssObjects[name]
	if !ok {
		return fmt.Errorf("DeleteObject object not exists %s", name)
	}
	delete(a.mc.OssObjects, name)
	return nil
}
