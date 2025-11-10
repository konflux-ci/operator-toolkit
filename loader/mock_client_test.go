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

package loader

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("MockClient", func() {

	Context("when mocking Create operations", func() {
		It("should return the mocked error when Create is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-configmap",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewBadRequest("mocked create error")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationCreate,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.Create(mockContext, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(expectedErr))
		})

		It("should use real client when no mock is configured", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "real-configmap",
					Namespace: "default",
				},
			}

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{},
			)

			err := mockClient.Create(mockContext, configMap)
			Expect(err).ToNot(HaveOccurred())

			// Verify it was actually created
			retrievedCM := &v1.ConfigMap{}
			err = k8sClient.Get(ctx, client.ObjectKey{Name: "real-configmap", Namespace: "default"}, retrievedCM)
			Expect(err).ToNot(HaveOccurred())

			// Cleanup
			Expect(k8sClient.Delete(ctx, configMap)).To(Succeed())
		})
	})

	Context("when mocking Update operations", func() {
		It("should return the mocked error when Update is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "update-test-cm",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewConflict(schema.GroupResource{}, "update-test-cm", nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationUpdate,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.Update(mockContext, configMap)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsConflict(err)).To(BeTrue())
		})
	})

	Context("when mocking Delete operations", func() {
		It("should return the mocked error when Delete is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "delete-test-cm",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewNotFound(schema.GroupResource{}, "delete-test-cm")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationDelete,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.Delete(mockContext, configMap)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})
	})

	Context("when mocking Patch operations", func() {
		It("should return the mocked error when Patch is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "patch-test-cm",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewInternalError(nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationPatch,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			patch := client.MergeFrom(configMap.DeepCopy())
			err := mockClient.Patch(mockContext, configMap, patch)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsInternalError(err)).To(BeTrue())
		})
	})

	Context("when mocking Get operations", func() {
		It("should return the mocked error when Get is called", func() {
			configMap := &v1.ConfigMap{}
			expectedErr := errors.NewNotFound(schema.GroupResource{}, "not-found-cm")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationGet,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.Get(mockContext, client.ObjectKey{Name: "not-found-cm", Namespace: "default"}, configMap)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should return the mocked result when Get is called", func() {
			expectedCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "mocked-cm",
					Namespace: "default",
					Labels: map[string]string{
						"test": "mocked",
					},
				},
			}

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationGet,
						ObjectType: &v1.ConfigMap{},
						Result:     expectedCM,
					},
				},
			)

			retrievedCM := &v1.ConfigMap{}
			err := mockClient.Get(mockContext, client.ObjectKey{Name: "mocked-cm", Namespace: "default"}, retrievedCM)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrievedCM.Name).To(Equal("mocked-cm"))
			Expect(retrievedCM.Labels).To(HaveKeyWithValue("test", "mocked"))
		})
	})

	Context("when mocking Status().Update() operations", func() {
		It("should return the mocked error when Status().Update() is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "status-update-test",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewForbidden(schema.GroupResource{}, "status-update-test", nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:       OperationUpdate,
						ObjectType:      &v1.ConfigMap{},
						SubResourceName: "status",
						Err:             expectedErr,
					},
				},
			)

			err := mockClient.Status().Update(mockContext, configMap)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsForbidden(err)).To(BeTrue())
		})
	})

	Context("when mocking multiple operations", func() {
		It("should handle different mocks for different object types", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "cm-multi-test",
					Namespace: "default",
				},
			}

			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-multi-test",
					Namespace: "default",
				},
			}

			cmErr := errors.NewBadRequest("configmap error")
			secretErr := errors.NewUnauthorized("secret error")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationCreate,
						ObjectType: &v1.ConfigMap{},
						Err:        cmErr,
					},
					{
						Operation:  OperationCreate,
						ObjectType: &v1.Secret{},
						Err:        secretErr,
					},
				},
			)

			// Test ConfigMap mock
			err := mockClient.Create(mockContext, configMap)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(cmErr))

			// Test Secret mock
			err = mockClient.Create(mockContext, secret)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(secretErr))
		})
	})

	Context("when combining loader mocks and client mocks", func() {
		It("should support both legacy loader mocks and new client mocks", func() {
			var loaderKey ContextKey = 100

			loaderResource := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name: "loader-resource",
				},
			}

			createErr := errors.NewBadRequest("create failed")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{
					{
						ContextKey: loaderKey,
						Resource:   loaderResource,
					},
				},
				[]ClientCallMock{
					{
						Operation:  OperationCreate,
						ObjectType: &v1.ConfigMap{},
						Err:        createErr,
					},
				},
			)

			// Test loader mock still works
			resource, err := GetMockedResourceAndErrorFromContext(mockContext, loaderKey, &v1.ConfigMap{})
			Expect(err).ToNot(HaveOccurred())
			Expect(resource.Name).To(Equal("loader-resource"))

			// Test client mock works
			newCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-cm",
					Namespace: "default",
				},
			}
			err = mockClient.Create(mockContext, newCM)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(createErr))
		})
	})

	Context("when mock client delegates to real client", func() {
		It("should call real client when no matching mock exists", func() {
			// Create a real resource first
			realCM := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "real-resource",
					Namespace: "default",
				},
				Data: map[string]string{
					"key": "value",
				},
			}
			Expect(k8sClient.Create(ctx, realCM)).To(Succeed())

			// Mock only Create operations, not Get
			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationCreate,
						ObjectType: &v1.ConfigMap{},
						Err:        errors.NewBadRequest("create error"),
					},
				},
			)

			// Get should use real client since no mock for Get exists
			retrievedCM := &v1.ConfigMap{}
			err := mockClient.Get(mockContext, client.ObjectKey{Name: "real-resource", Namespace: "default"}, retrievedCM)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrievedCM.Data).To(HaveKeyWithValue("key", "value"))

			// Cleanup
			Expect(k8sClient.Delete(ctx, realCM)).To(Succeed())
		})
	})

	Context("when mocking List operations", func() {
		It("should return the mocked error when List is called", func() {
			configMapList := &v1.ConfigMapList{}
			expectedErr := errors.NewTimeoutError("list operation timeout", 5)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationList,
						ObjectType: &v1.ConfigMapList{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.List(mockContext, configMapList)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsTimeout(err)).To(BeTrue())
		})

		It("should use real client when no mock is configured", func() {
			// Create a configmap first
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "list-test-cm",
					Namespace: "default",
				},
			}
			Expect(k8sClient.Create(ctx, configMap)).To(Succeed())

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{},
			)

			configMapList := &v1.ConfigMapList{}
			err := mockClient.List(mockContext, configMapList)
			Expect(err).ToNot(HaveOccurred())

			// Cleanup
			Expect(k8sClient.Delete(ctx, configMap)).To(Succeed())
		})
	})

	Context("when mocking Status().Patch() operations", func() {
		It("should return the mocked error when Status().Patch() is called", func() {
			configMap := &v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "status-patch-test",
					Namespace: "default",
				},
			}

			expectedErr := errors.NewServiceUnavailable("status patch unavailable")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:       OperationPatch,
						ObjectType:      &v1.ConfigMap{},
						SubResourceName: "status",
						Err:             expectedErr,
					},
				},
			)

			patch := client.MergeFrom(configMap.DeepCopy())
			err := mockClient.Status().Patch(mockContext, configMap, patch)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsServiceUnavailable(err)).To(BeTrue())
		})
	})

	Context("when mocking DeleteAllOf operations", func() {
		It("should return the mocked error when DeleteAllOf is called", func() {
			expectedErr := errors.NewForbidden(schema.GroupResource{Group: "", Resource: "configmaps"}, "", nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:  OperationDeleteAllOf,
						ObjectType: &v1.ConfigMap{},
						Err:        expectedErr,
					},
				},
			)

			err := mockClient.DeleteAllOf(mockContext, &v1.ConfigMap{}, client.InNamespace("default"))
			Expect(err).To(HaveOccurred())
			Expect(errors.IsForbidden(err)).To(BeTrue())
		})
	})

	Context("when mocking SubResource operations", func() {
		It("should return the mocked error when SubResource().Get() is called", func() {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			}
			subResource := &v1.Pod{}
			expectedErr := errors.NewNotFound(schema.GroupResource{}, "test-pod")

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:       OperationGet,
						ObjectType:      &v1.Pod{},
						SubResourceName: "status",
						Err:             expectedErr,
					},
				},
			)

			err := mockClient.SubResource("status").Get(mockContext, pod, subResource)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsNotFound(err)).To(BeTrue())
		})

		It("should return the mocked error when SubResource().Update() is called", func() {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod-update",
					Namespace: "default",
				},
			}
			expectedErr := errors.NewConflict(schema.GroupResource{}, "test-pod-update", nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:       OperationUpdate,
						ObjectType:      &v1.Pod{},
						SubResourceName: "status",
						Err:             expectedErr,
					},
				},
			)

			err := mockClient.SubResource("status").Update(mockContext, pod)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsConflict(err)).To(BeTrue())
		})

		It("should return the mocked error when SubResource().Patch() is called", func() {
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod-patch",
					Namespace: "default",
				},
			}
			expectedErr := errors.NewInternalError(nil)

			mockContext, mockClient := GetMockedContextWithClient(
				ctx,
				k8sClient,
				[]MockData{},
				[]ClientCallMock{
					{
						Operation:       OperationPatch,
						ObjectType:      &v1.Pod{},
						SubResourceName: "status",
						Err:             expectedErr,
					},
				},
			)

			patch := client.MergeFrom(pod.DeepCopy())
			err := mockClient.SubResource("status").Patch(mockContext, pod, patch)
			Expect(err).To(HaveOccurred())
			Expect(errors.IsInternalError(err)).To(BeTrue())
		})
	})
})
