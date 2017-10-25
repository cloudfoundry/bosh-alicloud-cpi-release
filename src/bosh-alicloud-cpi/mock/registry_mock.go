package mock

import "bosh-alicloud-cpi/registry"

func NewRegistryMock() (registry.Client) {
	return RegistryClientMock{}
}

type RegistryClientMock struct {
}


func (c RegistryClientMock) Delete(instanceID string) error {
	return nil
}

func (c RegistryClientMock) Fetch(instanceID string) (registry.AgentSettings, error) {
	return registry.AgentSettings{}, nil
}

func (c RegistryClientMock) Update(instanceID string, agentSettings registry.AgentSettings) error {
	return nil
}