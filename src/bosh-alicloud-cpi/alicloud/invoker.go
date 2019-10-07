/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"strings"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Invoker struct {
	catchers []*Catcher
}

type Catcher struct {
	Reason           string
	RetryCount       int
	RetryWaitSeconds int
}

var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 60, 5}
var OperationConflictCatcher = Catcher{"OperationConflict", 60, 5}
var InternalErrorCatcher = Catcher{"InternalError", 60, 5}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(OperationConflictCatcher)
	i.AddCatcher(InternalErrorCatcher)
	return i
}

func (a *Invoker) AddCatcher(catcher Catcher) {
	a.catchers = append(a.catchers, &catcher)
}

func (a *Invoker) Run(f func() error) error {
	err := f()

	if err == nil {
		return nil
	}

	for _, catcher := range a.catchers {
		if strings.Contains(err.Error(), catcher.Reason) {
			catcher.RetryCount--

			if catcher.RetryCount <= 0 {
				return bosherr.WrapError(err, "over max retry")
			} else {
				time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
				return a.Run(f)
			}
		}
	}
	return err
}

func (a *Invoker) RunUntil(timeout time.Duration, interval time.Duration, f func() (bool, error)) (bool, error) {
	for {
		ok, err := f()

		if err != nil {
			return false, bosherr.WrapError(err, "RunUntil failed")
		}

		if ok {
			return true, nil
		}

		timeout -= interval
		if timeout < 0 {
			return false, nil
		}
		time.Sleep(time.Duration(interval))
	}
}
