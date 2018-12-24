package alicloud

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	AlicloudOssServiceTag = "AlicloudOssService"
)

type OssManager interface {
	CreateBucket(name string, options ...oss.Option) error
	DeleteBucket(name string) error
	GetBucket(name string) (bucket *oss.Bucket, err error)
	UploadFile(bucket oss.Bucket, objectKey, filePath string, partSize int64, options ...oss.Option) error
	DeleteObject(bucket oss.Bucket, name string) error
}

type OssManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

func NewOssManager(config Config, logger boshlog.Logger) OssManager {
	return OssManagerImpl{
		config: config,
		logger: logger,
		region: config.OpenApi.GetRegion(""),
	}
}

func (a OssManagerImpl) CreateBucket(name string, options ...oss.Option) error {
	client, err := a.config.NewOssClient(a.region)
	if err != nil {
		return err
	}
	a.logger.Debug(AlicloudOssServiceTag, "Creating Alicloud Oss '%s'", name)

	err = client.CreateBucket(name, options...)
	return err
}

func (a OssManagerImpl) DeleteBucket(name string) error {
	client, err := a.config.NewOssClient(a.region)
	if err != nil {
		return err
	}
	a.logger.Debug(AlicloudOssServiceTag, "Deleting Alicloud Oss '%s'", name)

	return client.DeleteBucket(name)
}

func (a OssManagerImpl) GetBucket(name string) (bucket *oss.Bucket, err error) {
	client, err := a.config.NewOssClient(a.region)
	if err != nil {
		return nil, err
	}
	a.logger.Debug(AlicloudOssServiceTag, "Geting Alicloud Oss '%s'", name)

	return client.Bucket(name)
}

func (a OssManagerImpl) UploadFile(
	bucket oss.Bucket, objectKey, filePath string, partSize int64, options ...oss.Option) error {
	a.logger.Debug(AlicloudOssServiceTag, "Upload file '%s' to bucket %s.", objectKey, bucket.BucketName)
	return bucket.UploadFile(objectKey, filePath, partSize, options...)
}

func (a OssManagerImpl) DeleteObject(bucket oss.Bucket, name string) error {
	a.logger.Debug(AlicloudOssServiceTag, "Deleting Alicloud Object '%s'", name)
	return bucket.DeleteObject(name)
}
