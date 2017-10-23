/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"strings"
	"os"
)

type Builder string

func NewBuilder(json string) (Builder) {
	return Builder(json)
}

func (b Builder) P(p string, v string) (Builder) {
	if !strings.HasPrefix(p, "${") {
		p = "${" + p + "}"
	}
	return Builder(strings.Replace(string(b), p, v, 1))
}

func (b Builder) Env(p string) (Builder) {
	v := os.Getenv(p)
	if v != "" {
		return b.P(p, v)
	} else {
		return b
	}
}

func (b Builder) ToBytes() []byte {
	return []byte(b)
}