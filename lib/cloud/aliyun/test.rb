$LOAD_PATH << '.'

require './cloud.rb'

options = {}
options["aliyun"] = {}
options["aliyun"]["RegionId"] = "cn-hangzhou"
options["aliyun"]["InstanceType"] = "ecs.s3.large"
options["aliyun"]["ImageId"] = "m-23g9tihvk"
options["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
options["aliyun"]["InternetChargeType"] = "PayByTraffic"
options["aliyun"]["InternetMaxBandwidthOut"] = "3"
options["aliyun"]["Password"] = "1qaz@WSX"
options["aliyun"]["AccessKeyId"] = "***REMOVED***"
options["aliyun"]["Secret"] = "***REMOVED***"

cloud = Cloud.new(options)
cloud.create_vm()