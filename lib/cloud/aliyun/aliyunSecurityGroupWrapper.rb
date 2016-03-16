module Bosh::Aliyun
  class AliyunSecurityGroupWrapper
    #必传参数：地域(RegionId)
    def AliyunSecurityGroupWrapper.describeSecurityGroups(parameters)
      #parameter check:地域(RegionId)、安全组所在的专有网络(VpcId)、标签 key(Tag.n.Key)、标签value(Tag.n.Value)、当前页码(PageNumber)、分页查询时设置的每页行数(PageSize)
      parameters["Action"]= "DescribeSecurityGroups";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：目标安全组ID(SecurityGroupId)、目标安全组所属 Region ID(RegionId)
    def AliyunSecurityGroupWrapper.describeSecurityGroupAttribute(parameters)
      #parameter check:目标安全组ID(SecurityGroupId)、目标安全组所属 Region ID(RegionId)、网络类型(NicType)、授权方向(Direction)
      parameters["Action"]= "DescribeSecurityGroupAttribute";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：目标安全组ID(SecurityGroupId)、目标安全组所属 Region ID(RegionId)
    def AliyunSecurityGroupWrapper.hasSecurityGroup(parameters)
      #parameter check:目标安全组ID(SecurityGroupId)、目标安全组所属 Region ID(RegionId)、网络类型(NicType)、授权方向(Direction)
      AliyunSecurityGroupWrapper.describeSecurityGroupAttribute(parameters);
      return true;
    end
  end
end