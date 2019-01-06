/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"fmt"
	"os"
	"strings"
)

type Builder string

func NewBuilder(json string) Builder {
	return Builder(json)
}

func (b Builder) P(p string, v string) Builder {
	if !strings.HasPrefix(p, "${") {
		p = "${" + p + "}"
	}
	return Builder(strings.Replace(string(b), p, v, 1))
}

func (b Builder) ApplyEnvs() (Builder, error) {
	s := os.ExpandEnv(string(b))
	env := GetWrappedString(s, "${", "}")
	if s != "" {
		return b, fmt.Errorf("need env %s", env)
	}
	return Builder(s), nil
}

func GetWrappedString(s string, prefix string, suffix string) string {
	i := strings.Index(s, prefix)
	if i < 0 {
		return ""
	}
	j := strings.Index(s, suffix)
	if j < 0 {
		return string(s[i:])
	} else {
		return string(s[i:j])
	}
}

func (b Builder) ToBytes() []byte {
	return []byte(b)
}
