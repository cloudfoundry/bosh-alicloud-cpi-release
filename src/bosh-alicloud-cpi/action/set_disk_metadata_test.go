/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("cpi:set_disk_metadata", func() {
	It("can run reboot vm", func() {
		instCid, _ := mockContext.NewInstance()
		diskCid, _ := mockContext.NewDisk(instCid)
		r :=  caller.Run(mock.NewBuilder(`{
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
		Expect(disk.DiskName).Should(Equal("consul.441e940e-2ffe-4208-993e-3e5f888e2b7e_0"))
	})

	It("instance name can be normalized", func() {
		s := normalizeName("dea_ng/4daa76e3-db6d-4550-b9b4-c504e7865f62")
		Expect(s).Should(Equal("dea_ng.4daa76e3-db6d-4550-b9b4-c504e7865f62"))

		s = normalizeName("a")
		Expect(s).Should(Equal("i_a"))

		s = normalizeName("11123")
		Expect(s).Should(Equal( "i_11123"))
	})
})


