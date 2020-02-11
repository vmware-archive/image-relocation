/*
 * Copyright (c) 2018-Present Pivotal Software, Inc. All rights reserved.
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

var _ = Describe("FlattenRepoPath", func() {
	var (
		name   image.Name
		mapped string
	)

	JustBeforeEach(func() {
		mappedImage, err := pathmapping.FlattenRepoPath("test.host/testuser", name)
		Expect(err).NotTo(HaveOccurred())
		mapped = mappedImage.String()

		// check that the mapped path is valid
		_, err = image.NewName(mapped)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when the image path has a single element", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should flatten the path correctly", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-9482d6a53a1789fb7304a4fe88362903"))
		})
	})

	Context("when the image path has more than a single element", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some/path")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should flatten the path correctly", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-some-path-3236c106420c1d0898246e1d2b6ba8b6"))
		})
	})

	Context("when the image has a tag", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some/path:v1")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should flatten the path correctly", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-some-path-3236c106420c1d0898246e1d2b6ba8b6"))
		})
	})

	Context("when the image has a digest", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some/path@sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should flatten the path correctly", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-some-path-3236c106420c1d0898246e1d2b6ba8b6"))
		})
	})

	Context("when the image has a digest and a tag", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/some/path:v1@sha256:1e725169f37aec55908694840bc808fb13ebf02cb1765df225437c56a796f870")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should flatten the path correctly", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-some-path-3236c106420c1d0898246e1d2b6ba8b6"))
		})
	})

	Context("when the image path is long and has many elements", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/axxxxxxxxx/bxxxxxxxxx/cxxxxxxxxx/dxxxxxxxxx/" +
				"exxxxxxxxx/fxxxxxxxxx/gxxxxxxxxx/hxxxxxxxxx/ixxxxxxxxx/jxxxxxxxxx/kxxxxxxxxx/lxxxxxxxxx/" +
				"mxxxxxxxxx/nxxxxxxxxx/oxxxxxxxxx/pxxxxxxxxx/qxxxxxxxxx/rxxxxxxxxx/sxxxxxxxxx/txxxxxxxxx")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should omit some portions of the path", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-axxxxxxxxx-bxxxxxxxxx-cxxxxxxxxx-dxxxxxxxxx-exxxxxxxxx-fxxxxxxxxx-gxxxxxxxxx-hxxxxxxxxx-ixxxxxxxxx-jxxxxxxxxx-kxxxxxxxxx-lxxxxxxxxx-mxxxxxxxxx-nxxxxxxxxx-oxxxxxxxxx-pxxxxxxxxx---txxxxxxxxx-b0d16e8b4d43f2ec842cfcc61989a966"))
		})
	})

	Context("when the image path is long and has few elements", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/axxxxxxxxxabxxxxxxxxxbcxxxxxxxxxcdxxxxxxxxxd" +
				"exxxxxxxxxefxxxxxxxxx/gxxxxxxxxxghxxxxxxxxxhixxxxxxxxxijxxxxxxxxxjkxxxxxxxxxklxxxxxxxxxl" +
				"mxxxxxxxxxmnxxxxxxxxxnoxxxxxxxxxopxxxxxxxxxpqxxxxxxxxxqrxxxxxxxxxrsxxxxxxxxx/txxxxxxxxx")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should omit some portions of the path", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-axxxxxxxxxabxxxxxxxxxbcxxxxxxxxxcdxxxxxxxxxdexxxxxxxxxefxxxxxxxxx---txxxxxxxxx-4817a2fce97ff7aae687d14e66328781"))
		})
	})

	Context("when the image path is long and has two elements", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/axxxxxxxxxabxxxxxxxxxbcxxxxxxxxxcdxxxxxxxxxd" +
				"exxxxxxxxxefxxxxxxxxxfgxxxxxxxxxghxxxxxxxxxhixxxxxxxxxijxxxxxxxxxjkxxxxxxxxxklxxxxxxxxxl" +
				"mxxxxxxxxxmnxxxxxxxxxnoxxxxxxxxxopxxxxxxxxxpqxxxxxxxxxqrxxxxxxxxxrsxxxxxxxxxstxxxxxxxxx")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should omit some portions of the path", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user-a363dc80420c33618202ba2828aec456"))
		})
	})

	Context("when the image path is long and has three elements", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/some-user/axxxxxxxxxabxxxxxxxxxbcxxxxxxxxxcdxxxxxxxxxd" +
				"exxxxxxxxxefxxxxxxxxxfgxxxxxxxxxghxxxxxxxxxhixxxxxxxxxijxxxxxxxxxjkxxxxxxxxxklxxxxxxxxxl" +
				"mxxxxxxxxxmnxxxxxxxxxnoxxxxxxxxxopxxxxxxxxxpqxxxxxxxxxqrxxxxxxxxxrsxxxxxxxxx/txxxxxxxxx")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should omit some portions of the path", func() {
			Expect(mapped).To(Equal("test.host/testuser/some-user---txxxxxxxxx-5134b4594954926468fe0bdc23f640ef"))
		})
	})

	Context("when the first element of the image path is long", func() {
		BeforeEach(func() {
			var err error
			name, err = image.NewName("some.registry.com/axxxxxxxxxabxxxxxxxxxbcxxxxxxxxxcdxxxxxxxxxd" +
				"exxxxxxxxxefxxxxxxxxxfgxxxxxxxxxghxxxxxxxxxhixxxxxxxxxijxxxxxxxxxjkxxxxxxxxxklxxxxxxxxxl" +
				"mxxxxxxxxxmnxxxxxxxxxnoxxxxxxxxxopxxxxxxxxxpqxxxxxxxxxqrxxxxxxxxxrsxxxxxxxxxstxxxxxxxxx/suffix")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should omit some portions of the path", func() {
			Expect(mapped).To(Equal("test.host/testuser/suffix-3fa7b8289050d7d4fe5d56f3098397a0"))
		})
	})
})
