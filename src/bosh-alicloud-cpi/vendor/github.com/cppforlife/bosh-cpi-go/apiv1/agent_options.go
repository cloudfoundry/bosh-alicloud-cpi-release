package apiv1

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type AgentOptions struct {
	Mbus      string   // e.g. "https://user:password@127.0.0.1:4321/agent"
	NTP       []string // e.g. ["0.us.pool.ntp.org"]. Ok to be empty
	Blobstore BlobstoreOptions
}

type BlobstoreOptions struct {
	Type    string `json:"provider"`
	Options map[string]interface{}
}

func (o AgentOptions) Validate() error {
	if o.Mbus == "" {
		return bosherr.Error("Must provide non-empty Mbus")
	}

	err := o.Blobstore.Validate()
	if err != nil {
		return bosherr.WrapError(err, "Validating Blobstore configuration")
	}

	return nil
}

func (o BlobstoreOptions) Validate() error {
	if o.Type == "" {
		return bosherr.Error("Must provide non-empty Type")
	}

	return nil
}
