package apiv1

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type AgentEnv interface {
	AttachSystemDisk(interface{})
	AttachEphemeralDisk(interface{})
	AttachPersistentDisk(DiskCID, interface{})
	DetachPersistentDisk(DiskCID)
	AsBytes() ([]byte, error)
	_final() // interface unimplementable from outside
}

type AgentEnvImpl struct {
	spec agentEnvSpec
}

var _ AgentEnv = &AgentEnvImpl{}

type agentEnvSpec struct {
	AgentID string `json:"agent_id"`

	VM VMSpec `json:"vm"`

	Mbus string   `json:"mbus"`
	NTP  []string `json:"ntp"`

	Blobstore BlobstoreSpec `json:"blobstore"`

	Networks NetworksSpec `json:"networks"`

	Disks DisksSpec `json:"disks"`

	Env EnvSpec `json:"env"`
}

type VMSpec struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type NetworksSpec map[string]NetworkSpec

type NetworkSpec struct {
	Type string `json:"type"`

	IP      string `json:"ip"`
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`

	DNS     []string `json:"dns"`
	Default []string `json:"default"`

	MAC string `json:"mac"`

	Preconfigured bool `json:"preconfigured"`
}

type DisksSpec struct {
	System     interface{}    `json:"system"`
	Ephemeral  interface{}    `json:"ephemeral"`
	Persistent PersistentSpec `json:"persistent"`
}

type PersistentSpec map[string]interface{}

type EnvSpec map[string]interface{}

type BlobstoreSpec struct {
	Provider string                 `json:"provider"`
	Options  map[string]interface{} `json:"options"`
}

func (ae *AgentEnvImpl) AttachSystemDisk(hint interface{}) {
	ae.spec.Disks.System = hint
}

func (ae *AgentEnvImpl) AttachEphemeralDisk(hint interface{}) {
	ae.spec.Disks.Ephemeral = hint
}

// todo better type for hint
func (ae *AgentEnvImpl) AttachPersistentDisk(cid DiskCID, hint interface{}) {
	spec := PersistentSpec{}

	if ae.spec.Disks.Persistent != nil {
		for k, v := range ae.spec.Disks.Persistent {
			spec[k] = v
		}
	}

	spec[cid.AsString()] = hint

	ae.spec.Disks.Persistent = spec
}

func (ae *AgentEnvImpl) DetachPersistentDisk(cid DiskCID) {
	spec := PersistentSpec{}

	if ae.spec.Disks.Persistent != nil {
		for k, v := range ae.spec.Disks.Persistent {
			spec[k] = v
		}
	}

	delete(spec, cid.AsString())

	ae.spec.Disks.Persistent = spec
}

func (ae AgentEnvImpl) AsBytes() ([]byte, error) {
	bytes, err := json.Marshal(ae.spec)
	if err != nil {
		return nil, bosherr.WrapError(err, "Marshalling agent env")
	}

	return bytes, nil
}

func (ae AgentEnvImpl) _final() {}
