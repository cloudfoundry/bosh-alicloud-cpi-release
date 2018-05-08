/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:set_disk_metadata", func() {
	It("can set disk metadata", func() {
		instCid, _ := mockContext.NewInstance()
		diskCid, _ := mockContext.NewDisk(instCid)
		r := caller.Run(mock.NewBuilder(`{
        	"method": "set_disk_metadata",
        	"arguments": [
                "${DISK_ID}",
                {
					"director": "my-bosh",
					"deployment": "cf",
					"instance_id": "${INST_ID}",
					"job": "consul",
					"instance_index": "0",
					"instance_name": "consul/441e940e-2ffe-4208-993e-3e5f888e2b7e",
					"attached_at": "2017-11-03T06:21:27Z"
                }
	        ],
    	    "context": {
                "director_uuid": "d5555ed6-7688-4aae-9dff-4c4507042f3d",
                "request_id": "cpi-201248"
        	}
		}`).P("INST_ID", instCid).P("DISK_ID", diskCid).ToBytes())
		Expect(r.GetError()).NotTo(HaveOccurred())

		disk, ok := mockContext.Disks[diskCid]
		Expect(ok).To(BeTrue())
		Expect(disk.DiskName).Should(Equal(""))
	})
})
