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
	"os"
	"bosh-alicloud-cpi/registry"
	"github.com/denverdino/aliyungo/ecs"
	"strings"
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

type OpenApi struct {
	RegionId        string   `json:"region_id"`
	AccessKeyId     string   `json:"access_key_id"`
	AccessKeySecret string   `json:"access_key_secret"`
}

type RegistryConfig struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
	Host     string `json:"address"`
	Port     int	`json:"port"`
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

type BlobstoreOptions struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func (c Config) Validate() error {
	// TODO check configuration validation
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
		return config, bosherr.WrapError(err, "Unmarshalling config contents")
	}

	config = ccs.Root.Properties
	config.OpenApi.ApplySystemEnv()
	err = config.Validate()
	if err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

func (a *OpenApi) ApplySystemEnv() {
	a.AccessKeyId = os.ExpandEnv(a.AccessKeyId)
	a.AccessKeySecret = os.ExpandEnv(a.AccessKeySecret)
}

func (a *RegistryConfig) ToInstanceUserData() string {
	endpoint := fmt.Sprintf("%s://%s:%s@%s:%d", a.Protocol, a.User, a.Password, a.Host, a.Port)
	json := fmt.Sprintf(`{"RegistryConfig":{"Endpoint":"%s"}}`, endpoint)
	return json
}

func (a *BlobstoreConfig) AsRegistrySettings() (registry.BlobstoreSettings) {
	return registry.BlobstoreSettings {
		Provider: a.Provider,
		Options: a.Options,
	}
}

func (c Config) NewEcsClient() (*ecs.Client) {
	return ecs.NewClient(c.OpenApi.AccessKeyId, c.OpenApi.AccessKeySecret)
}


func (c Config) GetHttpRegistryClient(logger boshlog.Logger) (registry.Client) {
	r := c.Registry

	if strings.Compare("", r.Host) == 0 {
		//
		// first start need skip this operation
		return nil
	}

	clientOptions := registry.ClientOptions {
		Protocol: r.Protocol,
		Host: r.Host,
		Port: r.Port,
		Username: r.User,
		Password: r.Password,
	}

	client := registry.NewHTTPClient(clientOptions, logger)
	return client
}