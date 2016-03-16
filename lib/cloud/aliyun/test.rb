$LOAD_PATH << '.'

require 'aliyunSecurityGroupWrapper'

parameters={};
parameters["AccessKeyId"]= "***REMOVED***";

parameters["RegionId"]= "us-west-1";


parameters["Secret"]= "***REMOVED***&";

AliyunSecurityGroupWrapper.describeSecurityGroups(parameters);