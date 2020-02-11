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

package pathmapping_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"
)

var _ = Describe("FlattenRepoPathPreserveTagDigest", func() {

	const expectedMappedPath = "test.host/testuser/some-user-some-path-f4cdc2223f0c472921033d606fa74a89"

	var (
		name                     image.Name
		mapped                   string
		mappedWithoutTagOrDigest string
		tag                      string
		digest                   string
	)

	JustBeforeEach(func() {
		result, err := pathmapping.FlattenRepoPathPreserveTagDigest("test.host/testuser", name)
		Expect(err).NotTo(HaveOccurred())
		mapped = result.String()
		tag = result.Tag()
		digest = result.Digest().String()
		mappedWithoutTagOrDigest = result.WithoutTagOrDigest().String()

		// check that the mapped path is valid
		_, err = image.NewName(mapped)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when the image has neither a tag nor a digest", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some-path")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should map the image correctly", func() {
			Expect(mapped).To(Equal(expectedMappedPath))
		})

		It("should not introduce a tag", func() {
			Expect(tag).To(BeEmpty())
		})

		It("should not introduce a digest", func() {
			Expect(digest).To(BeEmpty())
		})
	})

	Context("when the image has a tag", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some-path:v1")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should preserve the tag", func() {
			Expect(tag).To(Equal("v1"))
		})

		It("should not introduce a digest", func() {
			Expect(digest).To(BeEmpty())
		})

		It("should map the path without regard to the tag", func() {
			Expect(mappedWithoutTagOrDigest).To(Equal(expectedMappedPath))
		})
	})

	Context("when the image has a digest", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some-path@sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not introduce a tag", func() {
			Expect(tag).To(BeEmpty())
		})

		It("should preserve the digest", func() {
			Expect(digest).To(Equal("sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870"))
		})

		It("should map the path without regard to the digest", func() {
			Expect(mappedWithoutTagOrDigest).To(Equal(expectedMappedPath))
		})
	})

	Context("when the image has a digest and a tag", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some-path:v1@sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should preserve the tag", func() {
			Expect(tag).To(Equal("v1"))
		})

		It("should preserve the digest", func() {
			Expect(digest).To(Equal("sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870"))
		})

		It("should map the path without regard to the tag or the digest", func() {
			Expect(mappedWithoutTagOrDigest).To(Equal(expectedMappedPath))
		})
	})
})
