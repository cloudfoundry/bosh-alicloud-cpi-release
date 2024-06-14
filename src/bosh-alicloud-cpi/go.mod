module bosh-alicloud-cpi

go 1.20

require (
	github.com/alibabacloud-go/tea v1.2.1
	github.com/alibabacloud-go/tea-rpc v1.3.3
	github.com/alibabacloud-go/tea-utils/v2 v2.0.4
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.676
	github.com/aliyun/aliyun-oss-go-sdk v3.0.1+incompatible
	github.com/aliyun/credentials-go v1.2.7
	github.com/cloudfoundry/bosh-utils v0.0.407
	github.com/cppforlife/bosh-cpi-go v0.0.0-20180718174221-526823bbeafd
	github.com/google/uuid v1.4.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.27.10
)

require (
	github.com/alibabacloud-go/debug v0.0.0-20190504072949-9472017b5c68 // indirect
	github.com/alibabacloud-go/tea-rpc-utils v1.1.2 // indirect
	github.com/alibabacloud-go/tea-utils v1.3.5 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/charlievieth/fs v0.0.3 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/jmespath/go-jmespath v0.0.0-20180206201540-c2b33e8439af // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/alibabacloud-go/tea-rpc v1.3.3 => github.com/alibabacloud-go/tea-rpc v1.3.3

replace github.com/alibabacloud-go/tea-utils/v2 v2.0.4 => github.com/alibabacloud-go/tea-utils/v2 v2.0.4

replace github.com/aliyun/alibaba-cloud-sdk-go v1.62.676 => github.com/aliyun/alibaba-cloud-sdk-go v1.62.676

replace github.com/aliyun/credentials-go v1.2.7 => github.com/aliyun/credentials-go v1.2.7
