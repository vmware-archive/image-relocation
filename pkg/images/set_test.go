/*
 * Copyright (c) 2020-Present Pivotal Software, Inc. All rights reserved.
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

package images_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/images"
)

var _ = Describe("Image set", func() {
	const (
		refA = "example.com/u/a"
		refB = "example.com/u/b"
	)

	Describe("New", func() {
		var (
			ii  []string
			s   images.Set
			err error
		)

		JustBeforeEach(func() {
			s, err = images.New(ii...)
		})

		Context("when the input slice is empty", func() {
			BeforeEach(func() {
				ii = []string{}
			})

			It("should construct an empty set", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(images.Empty))
			})
		})

		Context("when the input slice is non-empty", func() {
			BeforeEach(func() {
				ii = []string{refA, refB}
			})

			It("should construct the correct set", func() {
				Expect(err).NotTo(HaveOccurred())

				v, err := images.New(refB, refA)
				Expect(err).NotTo(HaveOccurred())

				Expect(s).To(Equal(v))
			})
		})

		Context("when the input slice has an invalid image reference", func() {
			BeforeEach(func() {
				ii = []string{"::"}
			})

			It("should return a suitable error", func() {
				Expect(err).To(MatchError(`invalid image reference: "::"`))
			})
		})
	})

	Describe("Union", func() {
		It("should produce the correct result", func() {
			s, err := images.New(refA)
			Expect(err).NotTo(HaveOccurred())

			u, err := images.New(refB)
			Expect(err).NotTo(HaveOccurred())

			v, err := images.New(refA, refB)
			Expect(err).NotTo(HaveOccurred())

			Expect(s.Union(u)).To(Equal(v))
		})
	})

	Describe("Slice", func() {
		It("should produce the correct result", func() {
			s, err := images.New(refA, refB)
			Expect(err).NotTo(HaveOccurred())

			a, err := image.NewName(refA)
			Expect(err).NotTo(HaveOccurred())

			b, err := image.NewName(refB)
			Expect(err).NotTo(HaveOccurred())

			Expect(s.Slice()).To(ConsistOf(a, b))
		})
	})

	Describe("Strings", func() {
		It("should produce the correct result", func() {
			s, err := images.New(refB, refA)
			Expect(err).NotTo(HaveOccurred())

			Expect(s.Strings()).To(Equal([]string{refA, refB}))
		})
	})

	Describe("String", func() {
		It("should produce the correct result", func() {
			s, err := images.New(refB, refA)
			Expect(err).NotTo(HaveOccurred())

			Expect(s.String()).To(Equal(fmt.Sprintf("[%s, %s]", refA, refB)))
		})
	})

	Describe("Marshalling", func() {
		It("should marshall to the correct string and back again", func() {
			s, err := images.New(refA, refB)
			Expect(err).NotTo(HaveOccurred())

			sb, err := s.MarshalJSON()
			Expect(err).NotTo(HaveOccurred())

			Expect(string(sb)).To(Equal(fmt.Sprintf("[%q,%q]", refA, refB)))

			var u images.Set
			err = (&u).UnmarshalJSON(sb)
			Expect(err).NotTo(HaveOccurred())

			Expect(u).To(Equal(s))
		})

		It("should unmarshall a string containing 'null' correctly", func() {
			var u images.Set
			err := (&u).UnmarshalJSON([]byte("null"))
			Expect(err).NotTo(HaveOccurred())

			Expect(u).To(Equal(images.Empty))
		})
	})
})
