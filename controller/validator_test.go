/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Validator", func() {

	When("Validate is called", func() {
		It("should return successfully if no validation functions are passed", func() {
			result := Validate()
			Expect(result.Err).NotTo(HaveOccurred())
			Expect(result.Valid).To(BeTrue())
		})

		It("should return successfully if all validation functions succeed", func() {
			result := Validate([]ValidationFunction{
				func() *ValidationResult {
					return &ValidationResult{Valid: true}
				},
				func() *ValidationResult {
					return &ValidationResult{Valid: true}
				},
			}...)
			Expect(result.Err).NotTo(HaveOccurred())
			Expect(result.Valid).To(BeTrue())
		})

		It("should fail immediately if a validation function fails", func() {
			result := Validate([]ValidationFunction{
				func() *ValidationResult {
					return &ValidationResult{Err: fmt.Errorf("validation failed")}
				},
				func() *ValidationResult {
					return &ValidationResult{Valid: true}
				},
			}...)
			Expect(result.Err).To(HaveOccurred())
			Expect(result.Err.Error()).To(Equal("validation failed"))
			Expect(result.Valid).To(BeFalse())
		})
	})
})
