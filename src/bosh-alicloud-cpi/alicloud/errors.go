/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	// common
	NotFound       = "NotFound"
	WaitForTimeout = "WaitForTimeout"
	// ecs
	InstanceNotFound        = "Instance.Notfound"
	RamInstanceNotFound     = "Forbidden.InstanceNotFound"
	MessageInstanceNotFound = "instance is not found"
	//stemcell
	ImageIsImporting = "ImageIsImporting"
)

var EcsInstanceNotFound = []string{"Instance.Notfound", "InvalidInstanceId.NotFound"}
var ResourceNotFound = []string{"InvalidResourceId.NotFound"}

// An Error represents a custom error for Terraform failure response
type ProviderError struct {
	errorCode string
	message   string
}

func (e *ProviderError) Error() string {
	return fmt.Sprintf("[ERROR] Bosh ALicloud CPI Error: Code: %s Message: %s.", e.errorCode, e.message)
}

func (err *ProviderError) ErrorCode() string {
	return err.errorCode
}

func (err *ProviderError) Message() string {
	return err.message
}

func GetNotFoundErrorFromString(str string) error {
	return &ProviderError{
		errorCode: InstanceNotFound,
		message:   str,
	}
}

func GetTimeErrorFromString(str string) error {
	return &ProviderError{
		errorCode: WaitForTimeout,
		message:   str,
	}
}

func GetNotFoundMessage(product, id string) string {
	return fmt.Sprintf("The specified %s %s is not found.", product, id)
}

func GetTimeoutMessage(product, status string) string {
	return fmt.Sprintf("Waitting for %s %s is timeout.", product, status)
}

func NotFoundError(err error) bool {

	if e, ok := err.(*errors.ServerError); ok &&
		(e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound ||
			strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)) {
		return true
	}

	if e, ok := err.(*ProviderError); ok &&
		(e.ErrorCode() == InstanceNotFound || e.ErrorCode() == RamInstanceNotFound || e.ErrorCode() == NotFound ||
			strings.Contains(strings.ToLower(e.Message()), MessageInstanceNotFound)) {
		return true
	}

	return false
}

func IsExceptedErrors(err error, expectCodes []string) bool {
	for _, code := range expectCodes {
		if e, ok := err.(*errors.ServerError); ok && (e.ErrorCode() == code || strings.Contains(e.Message(), code)) {
			return true
		}

		if e, ok := err.(*ProviderError); ok && (e.ErrorCode() == code || strings.Contains(e.Message(), code)) {
			return true
		}

		if e, ok := err.(oss.ServiceError); ok && (e.Code == code || strings.Contains(e.Message, code)) {
			return true
		}
	}
	return false
}
