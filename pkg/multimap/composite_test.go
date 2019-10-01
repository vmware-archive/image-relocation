/*
 * Copyright (c) 2019-Present Pivotal Software, Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package multimap_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/multimap"
)

var _ = Describe("Composite", func() {
	const (
		testInput       = "input"
		testOutput      = "output"
		testOutputOther = "another output"
		testInput2      = "input 2"
		testOutput2     = "output 2"
		testInput3      = "input 3"
		testOutput3     = "output 3"
	)
	var (
		stopCh  chan struct{}
		stopped bool
		comp    multimap.Composite
	)

	BeforeEach(func() {
		stopped = false
		stopCh = make(chan struct{})
		comp = multimap.New(stopCh)
	})

	AfterEach(func() {
		if !stopped {
			close(stopCh)
		}
	})

	It("should act as the identity function by default", func() {
		Expect(comp.Map(testInput)).To(Equal(testInput))
	})

	Context("when a mapping is added", func() {
		var mapping map[string]string
		BeforeEach(func() {
			mapping = map[string]string{testInput: testOutput, testInput2: testOutput2}
			Expect(comp.Add("test", mapping)).To(Succeed())
		})

		It("should apply the mapping", func() {
			Expect(comp.Map(testInput)).To(Equal(testOutput))
			Expect(comp.Map(testInput2)).To(Equal(testOutput2))
		})

		Context("when the mapping is deleted", func() {
			BeforeEach(func() {
				Expect(comp.Delete("test")).To(Succeed())
			})

			It("should no longer apply the mapping", func() {
				Expect(comp.Map(testInput)).To(Equal(testInput))
				Expect(comp.Map(testInput2)).To(Equal(testInput2))
			})
		})

		Context("when the mapping which was passed in is mutated", func() {
			BeforeEach(func() {
				mapping[testInput3] = testOutput3

				// Force composite to be recalculated
				comp.Add("test2", map[string]string{testInput: testOutput})
			})

			It("should not affect the composite mapping", func() {
				Expect(comp.Map(testInput3)).To(Equal(testInput3))
			})
		})

		Context("when another consistent mapping is added", func() {
			var err error
			BeforeEach(func() {
				err = comp.Add("test2", map[string]string{testInput: testOutput, testInput3: testOutput3})
			})

			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should apply the composed mapping", func() {
				Expect(comp.Map(testInput)).To(Equal(testOutput))
				Expect(comp.Map(testInput2)).To(Equal(testOutput2))
				Expect(comp.Map(testInput3)).To(Equal(testOutput3))
			})

			Context("when an inconsistent mapping is added in place of an existing mapping", func() {
				var err error
				BeforeEach(func() {
					err = comp.Add("test2", map[string]string{testInput: testOutputOther, testInput3: testOutput3})
				})

				It("should return an error", func() {
					Expect(err).To(MatchError("inconsistent mapping: test maps \"input\" to \"output\" but given mapping maps \"input\" to \"another output\""))
				})

				It("should retain the existing mapping", func() {
					Expect(comp.Delete("test")).To(Succeed())
				})

				It("should not leave behind the new mapping", func() {
					Expect(comp.Delete("test2")).To(MatchError("mapping not found: test2"))
				})

				It("should apply the composed mapping", func() {
					Expect(comp.Map(testInput)).To(Equal(testOutput))
					Expect(comp.Map(testInput2)).To(Equal(testOutput2))
					Expect(comp.Map(testInput3)).To(Equal(testInput3)) // testInput3 no longer in domain of composite mapping
				})
			})
		})

		Context("when an inconsistent new mapping is added", func() {
			var err error
			BeforeEach(func() {
				err = comp.Add("test2", map[string]string{testInput: testOutputOther})
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("inconsistent mapping: test maps \"input\" to \"output\" but given mapping maps \"input\" to \"another output\""))
			})

			It("should retain the existing mapping", func() {
				Expect(comp.Delete("test")).To(Succeed())
			})

			It("should not leave behind the new mapping", func() {
				Expect(comp.Delete("test2")).To(MatchError("mapping not found: test2"))
			})
		})
	})

	It("should return an error when deleting a non-existent mapping", func() {
		Expect(comp.Delete("test")).To(MatchError("mapping not found: test"))
	})

	Context("when the stop channel is closed", func() {
	    BeforeEach(func() {
	        close(stopCh)
	        stopped = true
	    })

	    It("should eventually panic on Add but function normally until then", func() {
			Eventually(func() bool {return panics(func() {
				Expect(comp.Add("test", map[string]string{})).To(Succeed())
			})}).Should(BeTrue())
	    })

	    It("should eventually panic on Delete but function normally until then", func() {
			Eventually(func() bool {return panics(func() {
				Expect(comp.Delete("test")).NotTo(Succeed())
			})}).Should(BeTrue())
	    })

	    It("should eventually panic on Map but function normally until then", func() {
			Eventually(func() bool {return panics(func() {
				Expect(comp.Map(testInput)).To(Equal(testInput))
			})}).Should(BeTrue())
		})
	})
})

func panics(f func()) (panicked bool) {
	defer func() {
		r := recover()
		panicked = r != nil
	}()
	f()
	return false
}
