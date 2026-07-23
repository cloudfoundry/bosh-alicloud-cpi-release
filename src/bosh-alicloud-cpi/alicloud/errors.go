/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"fmt"
	"strings"

	"github.com/alibabacloud-go/tea/tea"

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

// AliCloud ModifyDiskSpec refusal codes: target category not reachable in-place.
const (
	InvalidDiskCategoryNotSupported = "InvalidDiskCategory.NotSupported"
	OperationDeniedDiskCategory     = "OperationDenied.DiskCategoryNotSupport"
)

var DiskCategoryUnsupportedCodes = []string{
	InvalidDiskCategoryNotSupported,
	OperationDeniedDiskCategory,
}

// NotSupportedError signals the director that the disk mutation cannot be done
// in-place by AliCloud, but MAY be recoverable via snapshot-and-recreate.
type NotSupportedError struct {
	message string
}

func (e NotSupportedError) Error() string { return e.message }
func (e NotSupportedError) Type() string  { return "Bosh::Clouds::NotSupported" }

func NewNotSupportedError(format string, args ...interface{}) NotSupportedError {
	return NotSupportedError{message: fmt.Sprintf(format, args...)}
}

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

// NewProviderError constructs a *ProviderError with a caller-supplied code and
// message. Used by tests/mocks to synthesize a specific AliCloud error code.
func NewProviderError(errorCode, message string) *ProviderError {
	return &ProviderError{errorCode: errorCode, message: message}
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

		if e, ok := err.(*tea.SDKError); ok {
			for _, code := range expectCodes {
				// The second statement aims to match the tea sdk history bug
				if *e.Code == code || strings.HasPrefix(code, *e.Code) || strings.Contains(*e.Data, code) {
					return true
				}
			}
			return false
		}
	}
	return false
}
