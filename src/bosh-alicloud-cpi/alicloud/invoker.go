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

var ClientErrorCatcher = Catcher{"AliyunGoClientFailure", 5, 3}
var ServiceBusyCatcher = Catcher{"ServiceUnavailable", 5, 3}

func NewInvoker() Invoker {
	i := Invoker{}
	i.AddCatcher(ClientErrorCatcher)
	i.AddCatcher(ServiceBusyCatcher)
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
