package alicloud

type NetworkManager struct {
	runner Runner
}

func NewNetworkManager(runner Runner) NetworkManager {
	return NetworkManager{runner}
}
