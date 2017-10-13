package integration

import (
	. "github.com/onsi/gomega"
)

func assertSucceeds(request string) {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	Expect(response.Error).To(BeNil())
}

func assertSucceedsWithResult(request string) interface{} {
	response, err := execCPI(request)
	Expect(err).ToNot(HaveOccurred())
	print(response.Error.Message)
	Expect(response.Error).To(BeNil())

	Expect(response.Result).ToNot(BeNil())
	return response.Result
}
