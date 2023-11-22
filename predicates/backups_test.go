/*
Copyright 2023 Red Hat Inc.

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

package predicates

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/event"
)

var _ = Describe("Backups predicate", Ordered, func() {
	When("when IgnoreBackups predicate is used", func() {
		It("should process create events for objects without backup labels", func() {
			contextEvent := event.CreateEvent{Object: &corev1.Pod{}}
			Expect(IgnoreBackups{}.Create(contextEvent)).To(BeTrue())
		})

		It("should ignore create events for objects with backup labels", func() {
			contextEvent := event.CreateEvent{
				Object: &corev1.Pod{
					ObjectMeta: v1.ObjectMeta{
						Labels: map[string]string{"velero.io/backup-name": "foo"},
					},
				},
			}
			Expect(IgnoreBackups{}.Create(contextEvent)).To(BeFalse())
		})

		It("should process delete events", func() {
			contextEvent := event.DeleteEvent{Object: &corev1.Pod{}}
			Expect(IgnoreBackups{}.Delete(contextEvent)).To(BeTrue())
		})

		It("should process generic events", func() {
			contextEvent := event.GenericEvent{Object: &corev1.Pod{}}
			Expect(IgnoreBackups{}.Generic(contextEvent)).To(BeTrue())
		})

		It("should process update events", func() {
			contextEvent := event.CreateEvent{Object: &corev1.Pod{}}
			Expect(IgnoreBackups{}.Create(contextEvent)).To(BeTrue())
		})
	})
})
