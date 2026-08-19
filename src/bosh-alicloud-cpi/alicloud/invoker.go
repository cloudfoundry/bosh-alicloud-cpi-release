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

// A momentary network failure on the way to the API endpoint carries no API
// error code, so these match the text the net package produces instead. The
// budget is short because these clear in seconds or not at all.
var ConnectionResetCatcher = Catcher{"connection reset by peer", 6, 5}
var ConnectionTimeoutCatcher = Catcher{"i/o timeout", 6, 5}
var ConnectionRefusedCatcher = Catcher{"connection refused", 6, 5}
var DNSResolutionCatcher = Catcher{"no such host", 6, 5}
var TLSHandshakeTimeoutCatcher = Catcher{"TLS handshake timeout", 6, 5}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ServiceBusyCatcher)
	i.AddCatcher(OperationConflictCatcher)
	i.AddCatcher(InternalErrorCatcher)
	return i
}

// NewDescribeInvoker returns an Invoker that also retries a transient network
// failure reaching the API endpoint.
//
// Only describe calls get this. They are idempotent, so a retry cannot duplicate
// work; retrying a create or a delete whose response was merely lost could, so
// those keep using NewInvoker.
//
// The reason this exists: the status waits in ChangeInstanceStatus and
// ChangeDiskStatus end as soon as their status read returns an error, and the
// pipeline lost whole deploys to a single blip on that read -- DescribeInstances
// hitting `connection reset by peer`, and DescribeDisks failing to resolve the
// endpoint -- while the wait still had minutes of budget left and would have
// succeeded on its next tick. Absorbing the blip here keeps that timeout the one
// thing that bounds the wait. A network that is genuinely down still fails, only
// seconds later.
func NewDescribeInvoker() Invoker {
	i := NewInvoker()
	i.AddCatcher(ConnectionResetCatcher)
	i.AddCatcher(ConnectionTimeoutCatcher)
	i.AddCatcher(ConnectionRefusedCatcher)
	i.AddCatcher(DNSResolutionCatcher)
	i.AddCatcher(TLSHandshakeTimeoutCatcher)
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

	if catcher := a.catcherFor(err); catcher != nil {
		catcher.RetryCount--

		if catcher.RetryCount <= 0 {
			return bosherr.WrapError(err, "over max retry")
		} else {
			time.Sleep(time.Duration(catcher.RetryWaitSeconds) * time.Second)
			return a.Run(f)
		}
	}
	return err
}

// catcherFor returns the catcher that claims err, or nil when no catcher does
// and the error is therefore returned to the caller as-is.
func (a *Invoker) catcherFor(err error) *Catcher {
	for _, catcher := range a.catchers {
		if strings.Contains(err.Error(), catcher.Reason) {
			return catcher
		}
	}
	return nil
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
