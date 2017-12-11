package alicloud

import (
	"bosh-alicloud-cpi/registry"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"encoding/base64"
	"encoding/json"
)

type RegistryManager struct {
	config Config
	logger boshlog.Logger
}

func NewRegistryManager(c Config, logger boshlog.Logger) (RegistryManager) {
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
	client := a.config.NewEcsClient()
	args := ecs.DescribeUserdataArgs{
		RegionId: a.config.OpenApi.GetRegion(),
		InstanceId: instanceID,
	}

	var settings registry.AgentSettings

	r, err := client.DescribeUserdata(&args)
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
	a.log("Update/" + instanceID, err, data, "")
	return err
}

func (a RegistryManager) updateUserData(instanceID string, data []byte) (error) {
	client := a.config.NewEcsClient()

	var args ecs.ModifyInstanceAttributeArgs
	args.InstanceId = instanceID
	args.UserData = base64.StdEncoding.EncodeToString(data)

	invoker := NewInvoker()
	return invoker.Run(func() (error){
		return client.ModifyInstanceAttribute(&args)
	})
}