package mock

import (
	"bosh-alicloud-cpi/registry"
	"fmt"
)

func NewRegistryMock() (registry.Client) {
	r := RegistryClientMock{}
	r.registries = make(map[string]registry.AgentSettings)
	return r
}

type RegistryClientMock struct {
	registries map[string]registry.AgentSettings
}


func (c RegistryClientMock) Delete(instanceID string) error {
	delete(c.registries, instanceID)
	return nil
}

func (c RegistryClientMock) Fetch(instanceID string) (registry.AgentSettings, error) {
	r, ok := c.registries[instanceID]
	if ok {
		return r, nil
	} else {
		return registry.AgentSettings{}, fmt.Errorf("reigstry not found instance=%s", instanceID)
	}

	return registry.AgentSettings{}, nil
}

func (c RegistryClientMock) Update(instanceID string, agentSettings registry.AgentSettings) error {
	c.registries[instanceID] = agentSettings
	return nil
}