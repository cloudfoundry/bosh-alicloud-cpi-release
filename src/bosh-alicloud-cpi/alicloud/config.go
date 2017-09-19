package alicloud

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
	"os"
	"fmt"
	"encoding/base64"
	"bytes"
	"strings"
)

type Config struct {
	OpenApi OpenApi `json:"alicloud"`
	Actions Actions `json:"actions"`
}

type OpenApi struct {
	RegionId string			`json:"region_id"`
	AccessKeyId string		`json:"access_key_id"`
	AccessKeySecret string	`json:"access_key_secret"`
	Regions[] Region		`json:regions`
}

type Region struct {
	Name string				`json:"name"`
	ImageId string			`json:"image_id"`
}

type Actions struct {
	Agent Agent 			`json:"agent"`
	Registry Registry 		`json:"blobstore"`
}

type Agent struct {
	Mbus string 			`json:"mbus"`
	Blobstore Blobstore		`json:"blobstore"`
}

type Blobstore struct {
	Provider string			`json:"provider"`
	Options BlobstoreOptions			`json:"options"`

}

type BlobstoreOptions struct {
	Endpoint string			`json:"endpoint"`
	User string				`json:"agent"`
	Password string 		`json:"agent-password"`
}

type Registry struct {
	User string				`json:"user"`
	Password string			`json:"password"`
	Protocol string			`json:"protocol"`
	Host string				`json:"host"`
	Port string 			`json:"port"`
}

func (c Config) Validate() (error) {
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
	var config Config

	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config contents")
	}

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

func (a *Registry) ToInstanceUserData() (string) {
	endpoint := fmt.Sprintf("%s://%s:%s@%s:%s", a.Protocol, a.User, a.Password, a.Host, a.Port)
	json := fmt.Sprintf(`{"Registry":{"Endpoint":"%s"}}`, endpoint)
	return base64.StdEncoding.EncodeToString([]byte(json))
}