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

package metadata

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Metadata", func() {

	When("AddAnnotations is called", func() {
		It("should add the annotations to the object", func() {
			annotations := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{}

			Expect(pod.Annotations).To(HaveLen(0))
			Expect(AddAnnotations(pod, annotations)).To(Succeed())
			Expect(pod.Annotations).To(Equal(annotations))
		})

		It("should add the annotations to the existing ones", func() {
			annotations := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Annotations).To(HaveLen(1))
			Expect(AddAnnotations(pod, annotations)).To(Succeed())
			Expect(pod.Annotations).To(HaveLen(3))
			for key, value := range annotations {
				Expect(pod.Annotations).To(HaveKeyWithValue(key, value))
			}
		})

		It("should error if the object is nil", func() {
			err := AddAnnotations(nil, map[string]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("AddLabels is called", func() {
		It("should add the labels to the object", func() {
			labels := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{}

			Expect(pod.Labels).To(HaveLen(0))
			Expect(AddLabels(pod, labels)).To(Succeed())
			Expect(pod.Labels).To(Equal(labels))
		})

		It("should add the labels to the existing ones", func() {
			labels := map[string]string{"foo": "bar", "baz": "qux"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Labels).To(HaveLen(1))
			Expect(AddLabels(pod, labels)).To(Succeed())
			Expect(pod.Labels).To(HaveLen(3))
			for key, value := range labels {
				Expect(pod.Labels).To(HaveKeyWithValue(key, value))
			}
		})

		It("should error if the object is nil", func() {
			err := AddLabels(nil, map[string]string{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("CopyAnnotationsByPrefix is called", func() {
		It("should return all the annotations matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should copy all the annotations matching the given prefix when Annotations in destination is empty", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}

			Expect(CopyAnnotationsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
			}))
		})

		It("should copy nothing to the destination when annotations in source is nil ", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when source is nil", func() {
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsByPrefix(nil, podDest, "fo")).NotTo(BeNil())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when destination is nil", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsByPrefix(podSource, nil, "fo")).NotTo(BeNil())
		})
	})

	When("CopyAnnotationsWithPrefixReplacement is called", func() {
		It("should return all the annotations matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should copy all the annotations matching the given prefix when Annotations in destination is empty", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}

			Expect(CopyAnnotationsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao": Equal("bar"),
			}))
		})

		It("should copy nothing to the destination when annotations in source is nil ", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when source is nil", func() {
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsWithPrefixReplacement(nil, podDest, "fo", "ba")).NotTo(BeNil())
			Expect(podDest.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when destination is nil", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyAnnotationsWithPrefixReplacement(podSource, nil, "fo", "ba")).NotTo(BeNil())
		})
	})

	When("CopyLabelsByPrefix is called", func() {
		It("should return all the labels matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should copy all the labels matching the given prefix when Labels in destination is empty", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}

			Expect(CopyLabelsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
			}))
		})

		It("should copy nothing to the destination when Labels in source is nil ", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsByPrefix(podSource, podDest, "fo")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when source is nil", func() {
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsByPrefix(nil, podDest, "fo")).NotTo(BeNil())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when destination is nil", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsByPrefix(podSource, nil, "fo")).NotTo(BeNil())
		})
	})

	When("CopyLabelsWithPrefixReplacement is called", func() {
		It("should return all the labels matching the given prefix", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should copy all the labels matching the given prefix when Labels in destination is empty", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "bar", "baz": "qux"},
				},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}

			Expect(CopyLabelsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao": Equal("bar"),
			}))
		})

		It("should copy nothing to the destination when Labels in source is nil ", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{},
			}
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsWithPrefixReplacement(podSource, podDest, "fo", "ba")).To(Succeed())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when source is nil", func() {
			podDest := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsWithPrefixReplacement(nil, podDest, "fo", "ba")).NotTo(BeNil())
			Expect(podDest.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"quux": Equal("corge"),
			}))
		})

		It("should return error when destination is nil", func() {
			podSource := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(CopyLabelsWithPrefixReplacement(podSource, nil, "fo", "ba")).NotTo(BeNil())
		})
	})

	When("DeleteAnnotation is called", func() {
		It("should error if the object is nil", func() {
			err := DeleteAnnotation(nil, "foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

		It("should not fail when deleting an annotation from an object with no annotations", func() {
			pod := &corev1.Pod{}
			Expect(pod.Annotations).To(HaveLen(0))
			Expect(DeleteAnnotation(pod, "foo")).To(Succeed())
		})

		It("should not fail when deleting non-existent annotation from an object", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				}}

			Expect(pod.Annotations).To(HaveLen(1))
			Expect(DeleteAnnotation(pod, "foo")).To(Succeed())
			Expect(pod.Annotations).To(HaveLen(1))
		})

		It("should remove the specified annotation", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Annotations).To(HaveLen(1))
			Expect(DeleteAnnotation(pod, "quux")).To(Succeed())
			Expect(pod.Annotations).To(HaveLen(0))
		})

		It("should remove the specified annotation and preserve the others", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"quux": "corge",
						"foo":  "bar",
					},
				},
			}

			Expect(pod.Annotations).To(HaveLen(2))
			Expect(DeleteAnnotation(pod, "quux")).To(Succeed())
			Expect(pod.Annotations).To(HaveLen(1))
			Expect(pod.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
			}))
		})
	})

	When("DeleteLabel is called", func() {
		It("should error if the object is nil", func() {
			err := DeleteLabel(nil, "foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

		It("should not fail when deleting a label from an object with no labels", func() {
			pod := &corev1.Pod{}
			Expect(pod.Labels).To(HaveLen(0))
			Expect(DeleteLabel(pod, "foo")).To(Succeed())
		})

		It("should not fail when deleting a non-existent label from an object", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				}}

			Expect(pod.Labels).To(HaveLen(1))
			Expect(DeleteLabel(pod, "foo")).To(Succeed())
			Expect(pod.Labels).To(HaveLen(1))
		})

		It("should remove the specified label", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Labels).To(HaveLen(1))
			Expect(DeleteLabel(pod, "quux")).To(Succeed())
			Expect(pod.Labels).To(HaveLen(0))
		})

		It("should remove the specified label and preserve the others", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{
						"quux": "corge",
						"foo":  "bar",
					},
				},
			}

			Expect(pod.Labels).To(HaveLen(2))
			Expect(DeleteLabel(pod, "quux")).To(Succeed())
			Expect(pod.Labels).To(HaveLen(1))
			Expect(pod.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
			}))
		})
	})

	When("HasAnnotation is called", func() {
		It("should return true when the annotation is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasAnnotation(pod, "bar")).To(BeTrue())
		})

		It("should return false when the annotation is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasAnnotation(pod, "nobar")).To(BeFalse())
		})
	})

	When("HasAnnotationWithValue is called", func() {
		It("should return true when the annotation with the given value is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasAnnotationWithValue(pod, "bar", "foo")).To(BeTrue())
		})

		It("should return false when the annotation with the given value is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasAnnotationWithValue(pod, "bar", "nofoo")).To(BeFalse())
		})
	})

	When("HasLabel is called", func() {
		It("should return true when the label is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasLabel(pod, "bar")).To(BeTrue())
		})

		It("should return false when the label is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasLabel(pod, "nobar")).To(BeFalse())
		})
	})

	When("HasLabelWithValue is called", func() {
		It("should return true when the label with the given value is found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasLabelWithValue(pod, "bar", "foo")).To(BeTrue())
		})

		It("should return false when the label with the given value is not found", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			Expect(HasLabelWithValue(pod, "bar", "nofoo")).To(BeFalse())
		})
	})

	When("GetAnnotationsWithPrefix is called", func() {
		It("should return all the annotations matching the given prefix", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			annotations, err := GetAnnotationsWithPrefix(pod, "ba")
			Expect(annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an empty map if there are no matching annotations", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			annotations, err := GetAnnotationsWithPrefix(pod, "nil")
			Expect(annotations).To(HaveLen(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should error if the object is nil", func() {
			_, err := GetAnnotationsWithPrefix(nil, "prefix")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

	})

	When("GetLabelsWithPrefix is called", func() {
		It("should return all the labels matching the given prefix", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			labels, err := GetLabelsWithPrefix(pod, "ba")
			Expect(labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return an empty map if there are no matching labels", func() {
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"bar": "foo", "baz": "qux"},
				},
			}

			labels, err := GetLabelsWithPrefix(pod, "nil")
			Expect(labels).To(HaveLen(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("should error if the object is nil", func() {
			_, err := GetLabelsWithPrefix(nil, "prefix")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})
	})

	When("SetAnnotation is called", func() {
		It("should error if the object is nil", func() {
			err := SetAnnotation(nil, "foo", "bar")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

		It("should add the annotation to the object", func() {
			annotations := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{}

			Expect(pod.Annotations).To(HaveLen(0))

			for key, value := range annotations {
				Expect(SetAnnotation(pod, key, value)).To(Succeed())
			}

			Expect(pod.Annotations).To(Equal(annotations))

		})

		It("should add the annotation to the existing ones", func() {
			annotations := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Annotations).To(HaveLen(1))

			for key, value := range annotations {
				Expect(SetAnnotation(pod, key, value)).To(Succeed())
			}

			Expect(pod.Annotations).To(HaveLen(2))
			Expect(pod.Annotations).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should update the existing annotation", func() {
			annotations := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{"foo": "corge"},
				},
			}

			Expect(pod.Annotations).To(HaveLen(1))

			for key, value := range annotations {
				Expect(SetAnnotation(pod, key, value)).To(Succeed())
			}

			Expect(pod.Annotations).To(HaveLen(1))
			Expect(pod.Annotations).To(Equal(annotations))

		})
	})

	When("SetLabel is called", func() {
		It("should error if the object is nil", func() {
			err := SetLabel(nil, "foo", "bar")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("object cannot be nil"))
		})

		It("should add the label to the object", func() {
			labels := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{}

			Expect(pod.Labels).To(HaveLen(0))

			for key, value := range labels {
				Expect(SetLabel(pod, key, value)).To(Succeed())
			}

			Expect(pod.Labels).To(Equal(labels))
		})

		It("should add the label to the existing ones", func() {
			labels := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"quux": "corge"},
				},
			}

			Expect(pod.Labels).To(HaveLen(1))

			for key, value := range labels {
				Expect(SetLabel(pod, key, value)).To(Succeed())
			}

			Expect(pod.Labels).To(HaveLen(2))
			Expect(pod.Labels).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})

		It("should update the existing label", func() {
			labels := map[string]string{"foo": "bar"}
			pod := &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Labels: map[string]string{"foo": "corge"},
				},
			}

			Expect(pod.Labels).To(HaveLen(1))

			for key, value := range labels {
				Expect(SetLabel(pod, key, value)).To(Succeed())
			}

			Expect(pod.Labels).To(HaveLen(1))
			Expect(pod.Labels).To(Equal(labels))

		})
	})

	When("addEntries is called", func() {
		It("should merge both maps", func() {
			source := map[string]string{"foo": "bar", "baz": "qux"}
			destination := map[string]string{"quux": "corge"}

			Expect(destination).To(HaveLen(1))
			addEntries(source, destination)
			Expect(destination).To(HaveLen(3))
			for key, value := range source {
				Expect(destination).To(HaveKeyWithValue(key, value))
			}
		})
	})

	When("copyByPrefix is called", func() {
		It("should copy key/value pairs with the prefix to destination", func() {
			source := map[string]string{"foo": "bar", "foz": "qux"}
			destination := map[string]string{"quux": "corge"}
			copyByPrefix(source, destination, "foo")
			Expect(destination).To(HaveLen(2))
			Expect(destination).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo":  Equal("bar"),
				"quux": Equal("corge"),
			}))
		})
	})

	When("copyWithPrefixReplacement is called", func() {
		It("should copy key/value pairs with the prefix to destination with new prefix", func() {
			source := map[string]string{"foo": "bar", "foz": "qux"}
			destination := map[string]string{"quux": "corge"}
			copyWithPrefixReplacement(source, destination, "fo", "ba")
			Expect(destination).To(HaveLen(3))
			Expect(destination).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bao":  Equal("bar"),
				"baz":  Equal("qux"),
				"quux": Equal("corge"),
			}))
		})
	})

	When("filterByPrefix is called", func() {
		It("should return the matching elements", func() {
			entries := map[string]string{"bar": "foo", "baz": "qux"}

			matchingEntries := filterByPrefix(entries, "ba")
			Expect(matchingEntries).To(HaveLen(2))
			Expect(matchingEntries).To(gstruct.MatchAllKeys(gstruct.Keys{
				"bar": Equal("foo"),
				"baz": Equal("qux"),
			}))
		})
	})

	When("safeCopy is called", func() {
		It("should copy the given key into the destination map", func() {
			destination := map[string]string{}
			safeCopy(destination, "foo", "bar")
			Expect(destination).To(HaveLen(1))
			Expect(destination).To(HaveKeyWithValue("foo", "bar"))
		})

		It("should preserve the existing key/values", func() {
			destination := map[string]string{"baz": "qux"}
			safeCopy(destination, "foo", "bar")
			Expect(destination).To(HaveLen(2))
			Expect(destination).To(gstruct.MatchAllKeys(gstruct.Keys{
				"foo": Equal("bar"),
				"baz": Equal("qux"),
			}))
		})
	})
})
