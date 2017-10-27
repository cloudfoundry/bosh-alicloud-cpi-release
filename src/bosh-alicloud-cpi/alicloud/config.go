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
	UseForceStop			= true

	WaitTimeout  = time.Duration(120) * time.Second
	WaitInterval = time.Duration(3) * time.Second

	DefaultEipWaitSeconds	= 120
	DefaultSlbWeight = 100
)

type OpenApi struct {
	RegionId        string  `json:"region_id"`
	ZoneId			string	`json:"zone_id"`
	AccessKeyId     string  `json:"access_key_id"`
	AccessKeySecret string  `json:"access_key_secret"`
	Region			common.Region	`json:"-"`
}

type RegistryConfig struct {
	User     string			`json:"user"`
	Password string			`json:"password"`
	Protocol string			`json:"protocol"`
	Host     string			`json:"address"`
	Port     json.Number	`json:"port"`
}

type AgentConfig struct {
	Ntp       []string        `json:"ntp"`
	Mbus      string          `json:"mbus"`
	Blobstore BlobstoreConfig `json:"blobstore"`
}

type BlobstoreConfig struct {
	Provider string          		`json:"provider"`
	Options  map[string]interface{}	`json:"options"`
}

func (c Config) Validate() error {
	if c.OpenApi.RegionId == "" {
		return fmt.Errorf("region can't be empty")
	}
	c.OpenApi.Region = common.Region(c.OpenApi.RegionId)

	_, err := c.Registry.Port.Int64()
	if err != nil {
		return fmt.Errorf("bad registry.port %s", c.Registry.Port.String())
	}

	//TODO: validate more
	return nil
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
		return config, bosherr.WrapError(err, "Unmarshal config json failed")
	}

	config = ccs.Root.Properties

	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
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
	return registry.BlobstoreSettings {
		Provider: a.Provider,
		Options: a.Options,
	}
}

func (c Config) NewEcsClient() (*ecs.Client) {
	return ecs.NewClient(c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
}

func (c Config) NewSlbClient() (*slb.Client) {
	return slb.NewClient(c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
}


func (c Config) GetHttpRegistryClient(logger boshlog.Logger) (registry.Client) {
	r := c.Registry

	port, _ := r.Port.Int64()
	clientOptions := registry.ClientOptions {
		Protocol: r.Protocol,
		Host: r.Host,
		Port: int(port),
		Username: r.User,
		Password: r.Password,
	}

	client := registry.NewHTTPClient(clientOptions, logger)
	return client
}