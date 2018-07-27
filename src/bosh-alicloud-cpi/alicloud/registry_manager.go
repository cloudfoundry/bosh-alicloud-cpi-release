package alicloud

import (
	"bosh-alicloud-cpi/registry"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type RegistryManager struct {
	config Config
	logger boshlog.Logger
}

func NewRegistryManager(c Config, logger boshlog.Logger) RegistryManager {
	return RegistryManager{config: c, logger: logger}
}

func (a RegistryManager) log(action string, err error, args interface{}, result string) {
	s, _ := json.Marshal(args)
	if err != nil {
		a.logger.Error("RegistryManager", "%s failed args=%s err=%s", action, s, err)
	} else {
		a.logger.Info("RegistryManager", "%s done! args=%s result=%s", action, s, result)
	}
}

func (a RegistryManager) Delete(instanceID string) error {
	err := a.updateUserData(instanceID, []byte(""))
	a.log("Delete", err, instanceID, "")
	return err
}

func (a RegistryManager) Fetch(instanceID string) (registry.AgentSettings, error) {
	var settings registry.AgentSettings
	region, err := a.config.GetInstanceRegion(instanceID)
	if err != nil {
		return settings, err
	}
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return settings, err
	}

	args := ecs.CreateDescribeUserDataRequest()
	args.InstanceId = instanceID

	r, err := client.DescribeUserData(args)
	if err != nil {
		return settings, fmt.Errorf("OpenAPI::DescribeUserData Failed %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(r.UserData)
	if err != nil {
		return settings, fmt.Errorf("decode UserData %s failed %v", r.UserData, err)
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		return settings, fmt.Errorf("unmarshal UserData failed %v", err)
	}

	return settings, nil
}

func (a RegistryManager) Update(instanceID string, settings registry.AgentSettings) error {
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	err = a.updateUserData(instanceID, data)
	a.log("Update/"+instanceID, err, data, "")
	return err
}

func (a RegistryManager) updateUserData(instanceID string, data []byte) error {
	region, err := a.config.GetInstanceRegion(instanceID)
	if err != nil {
		return err
	}
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return err
	}

	args := ecs.CreateModifyInstanceAttributeRequest()
	args.InstanceId = instanceID
	args.UserData = base64.StdEncoding.EncodeToString(data)

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, err := client.ModifyInstanceAttribute(args)
		return err
	})
}
