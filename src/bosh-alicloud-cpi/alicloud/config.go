/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"encoding/json"
	"fmt"
	"bosh-alicloud-cpi/registry"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"time"
	"github.com/denverdino/aliyungo/slb"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
	"net"
)

type InnerType string

const (
	DefaultClassOSSInnerEndpoint = "oss-cn-hangzhou-internal"
	DefaultVpcOSSInnerEndpoint   = "oss-cn-hangzhou-internal"
	OSSSuffix                    = "oss-"

	InnerVpc       = InnerType("VPC")
	InnerClassic   = InnerType("CLASSIC")
	PingMethod     = "tcp"
	TimeoutSeconds = 5
)

type CloudConfigJson struct {
	Root CloudConfig `json:"cloud"`
}

type CloudConfig struct {
	Plugin     string `json:"plugin"`
	Properties Config `json:"properties"`
}

type Config struct {
	OpenApi  OpenApi        `json:"alicloud"`
	Registry RegistryConfig `json:"registry"`
	Agent    AgentConfig    `json:"agent"`
}

const (
	UseForceStop = true

	WaitTimeout  = time.Duration(180) * time.Second
	WaitInterval = time.Duration(5) * time.Second

	DefaultEipWaitSeconds = 120
	DefaultSlbWeight      = 100
)

type OpenApi struct {
	RegionId        string  `json:"region_id"`
	ZoneId			string	`json:"zone_id"`
	AccessEndpoint	string 	`json:"access_endpoint"`
	AccessKeyId     string  `json:"access_key_id"`
	AccessKeySecret string  `json:"access_key_secret"`
}

type RegistryConfig struct {
	User     string      `json:"user"`
	Password string      `json:"password"`
	Protocol string      `json:"protocol"`
	Host     string      `json:"host"`
	Port     json.Number `json:"port"`
}

type AgentConfig struct {
	Ntp       []string        `json:"ntp"`
	Mbus      string          `json:"mbus"`
	Blobstore BlobstoreConfig `json:"blobstore"`
}

type BlobstoreConfig struct {
	Provider string                 `json:"provider"`
	Options  map[string]interface{} `json:"options"`
}

func (c Config) Validate() error {
	if c.OpenApi.RegionId == "" {
		return fmt.Errorf("region can't be empty")
	}

	_, err := c.Registry.Port.Int64()
	if err != nil {
		return fmt.Errorf("bad registry.port %s", c.Registry.Port.String())
	}

	//TODO: validate more
	return nil
}

func (a OpenApi) GetRegion() (common.Region) {
	return common.Region(a.RegionId)
}

func (a RegistryConfig) IsEmpty() (bool) {
	if a.Host == "" {
		return true
	} else {
		return false
	}
}

func NewConfigFromFile(configFile string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	if configFile == "" {
		return config, bosherr.Errorf("Must provide a config file")
	}

	bytes, err := fs.ReadFile(configFile)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config file '%s'", configFile)
	}

	return NewConfigFromBytes(bytes)
}

func NewConfigFromBytes(bytes []byte) (Config, error) {
	var ccs CloudConfigJson
	var config Config

	err := json.Unmarshal(bytes, &ccs)
	if err != nil {
		return config, bosherr.WrapError(err, "unmarshal config json failed")
	}

	config = ccs.Root.Properties

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "validate config failed")
	}

	return config, nil
}

func (a RegistryConfig) ToInstanceUserData() string {
	endpoint := a.GetEndpoint()
	json := fmt.Sprintf(`{"registry":{"endpoint":"%s"}}`, endpoint)
	return json
}

func (a RegistryConfig) GetEndpoint() (string) {
	port, _ := a.Port.Int64()
	return fmt.Sprintf("%s://%s:%s@%s:%d", a.Protocol, a.User, a.Password, a.Host, port)
}

func (a BlobstoreConfig) AsRegistrySettings() (registry.BlobstoreSettings) {
	return registry.BlobstoreSettings{
		Provider: a.Provider,
		Options:  a.Options,
	}
}

func (c Config) NewEcsClient() (*ecs.Client) {
	// Obsoleted
	// ep := "https://ecs." + c.OpenApi.GetEndpoint()
	// return ecs.NewClientWithEndpoint(ep, c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
	return ecs.NewECSClient(c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret, common.Region(c.OpenApi.RegionId))
}

func (c Config) NewSlbClient() (*slb.Client) {
	// Obsoleted
	// ep := "https://slb." + c.OpenApi.GetEndpoint()
	// return slb.NewClientWithEndpoint(ep, c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
	return slb.NewSLBClient(c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret, common.Region(c.OpenApi.RegionId))
}

func (c Config) GetRegistryClient(logger boshlog.Logger) (registry.Client) {
	if !c.Registry.IsEmpty() {
		return c.GetHttpRegistryClient(logger)
	} else {
		return NewRegistryManager(c, logger)
	}
}

func (c Config) NewOssClient(inner bool) (*oss.Client) {
	ossClient, _ := oss.New(c.GetAvailableOSSEndPoint(inner), c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
	return ossClient
}

func (c Config) GetAvailableOSSEndPoint(inner bool) (string) {
	return "https://" + c.GetOSSEndPoint(inner) + ".aliyuncs.com"
}

func (c Config) GetOSSEndPoint(inner bool) (string) {
	timeOut := time.Duration(TimeoutSeconds) * time.Second
	ep := GetOSSEndPoint(string(c.OpenApi.GetRegion()), "")
	if !inner {
		return ep
	}

	ep = GetOSSEndPoint("", InnerVpc)
	if _, err := net.DialTimeout(PingMethod, ep, timeOut); err != nil {
		fmt.Printf("Ping oss inner vpc endpoint %s ok", ep)
		return ep
	}

	ep = GetOSSEndPoint("", InnerClassic)
	if _, err := net.DialTimeout(PingMethod, ep, timeOut); err != nil {
		fmt.Printf("Ping oss inner ecs endpoint %s ok", ep)
		return ep
	}

	ep = GetOSSEndPoint(string(c.OpenApi.GetRegion()), "")
	return ep
}

// types allows ["VPC", "CLASSIC"], then return inner endpoint
// otherwise return endpoint by region
func GetOSSEndPoint(region string, types InnerType) string {
	if types == InnerVpc {
		return DefaultVpcOSSInnerEndpoint
	}

	if types == InnerClassic {
		return DefaultClassOSSInnerEndpoint
	}

	if strings.HasPrefix(region, OSSSuffix) {
		return region
	}
	return OSSSuffix + region
}

func (c Config) GetHttpRegistryClient(logger boshlog.Logger) (registry.Client) {
	r := c.Registry

	port, _ := r.Port.Int64()
	clientOptions := registry.ClientOptions{
		Protocol: r.Protocol,
		Host:     r.Host,
		Port:     int(port),
		Username: r.User,
		Password: r.Password,
	}

	client := registry.NewHTTPClient(clientOptions, logger)
	return client
}
