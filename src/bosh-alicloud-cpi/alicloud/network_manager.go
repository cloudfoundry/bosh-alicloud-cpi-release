/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

type NetworkManager struct {
	runner Runner
}

func NewNetworkManager(runner Runner) NetworkManager {
	return NetworkManager{runner}
}
