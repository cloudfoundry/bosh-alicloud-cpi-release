#!/usr/bin/expect
spawn git pull https://xiaozhu36@github.com/xiaozhu36/bosh-alicloud-cpi-release.git concourse_ci_tmp
expect "Password for 'https://xiaozhu36@github.com': "
send "xiaozhu36@alibaba\r"
exit
