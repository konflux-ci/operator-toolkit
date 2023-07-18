/*
Copyright 2022.

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

package loader

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Loader", func() {

	When("GetMockedContext is called", func() {
		It("should return a new context with the given key and value", func() {
			var contextKey ContextKey = 0
			mockData := []MockData{
				{
					ContextKey: contextKey,
					Resource: &v1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "configMap",
							Namespace: "default",
						},
					},
				},
			}

			mockContext := GetMockedContext(ctx, mockData)
			returnedMockData := mockContext.Value(contextKey).(MockData)
			Expect(returnedMockData.ContextKey).To(Equal(mockData[0].ContextKey))
			Expect(returnedMockData.Resource).To(Equal(mockData[0].Resource))
		})
	})

	When("GetMockedResourceAndErrorFromContext is called", func() {
		var contextKey ContextKey = 0

		contextErr := errors.NewNotFound(schema.GroupResource{}, "")
		contextResource := &v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "configMap",
				Namespace: "default",
			},
		}

		It("returns the resource from the context", func() {
			mockContext := GetMockedContext(ctx, []MockData{
				{
					ContextKey: contextKey,
					Resource:   contextResource,
				},
			})
			resource, err := GetMockedResourceAndErrorFromContext(mockContext, contextKey, contextResource)
			Expect(err).To(BeNil())
			Expect(resource).To(Equal(contextResource))
		})

		It("returns the error from the context", func() {
			mockContext := GetMockedContext(ctx, []MockData{
				{
					ContextKey: contextKey,
					Err:        contextErr,
				},
			})
			resource, err := GetMockedResourceAndErrorFromContext(mockContext, contextKey, contextResource)
			Expect(err).To(Equal(contextErr))
			Expect(resource).To(BeNil())
		})

		It("returns the resource and the error from the context", func() {
			mockContext := GetMockedContext(ctx, []MockData{
				{
					ContextKey: contextKey,
					Resource:   contextResource,
					Err:        contextErr,
				},
			})
			resource, err := GetMockedResourceAndErrorFromContext(mockContext, contextKey, contextResource)
			Expect(err).To(Equal(contextErr))
			Expect(resource).To(Equal(contextResource))
		})

		It("should panic when the mocked data is not present", func() {
			Expect(func() {
				_, _ = GetMockedResourceAndErrorFromContext(ctx, contextKey, contextResource)
			}).To(Panic())
		})
	})

	When("GetObject is called", func() {
		It("returns the requested resource if it exists", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "config-map",
					Namespace: "default",
				},
			}
			Expect(k8sClient.Create(ctx, configMap)).To(Succeed())

			returnedObject := &v1.ConfigMap{}
			Expect(GetObject(configMap.Name, configMap.Namespace, k8sClient, ctx, returnedObject)).To(Succeed())
			Expect(configMap.ObjectMeta).To(Equal(returnedObject.ObjectMeta))

			Expect(k8sClient.Delete(ctx, configMap)).To(Succeed())
		})

		It("returns and error if the requested resource doesn't exist", func() {
			returnedObject := &v1.ConfigMap{}
			err := GetObject("non-existent", "non-existent", k8sClient, ctx, returnedObject)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})

})
