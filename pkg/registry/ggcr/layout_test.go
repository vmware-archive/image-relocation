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

package ggcr_test

import (
	"errors"

	"github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
	x "github.com/pivotal/image-relocation/pkg/registry/ggcr/ggcrfakes"
	"github.com/pivotal/image-relocation/pkg/registry/ggcrfakes"
	"github.com/pivotal/image-relocation/pkg/registry/registryfakes"
)

var _ = Describe("Layout", func() {
	var (
		layout             registry.Layout
		mockLayoutPath     *registryfakes.FakeLayoutPath
		mockImageIndex     *ggcrfakes.FakeImageIndex
		mockRegistryClient *x.FakeRegistryClient
		testError          error
	)

	BeforeEach(func() {
		mockLayoutPath = &registryfakes.FakeLayoutPath{}
		mockImageIndex = &ggcrfakes.FakeImageIndex{}
		mockRegistryClient = &x.FakeRegistryClient{}

		layout = ggcr.NewImageLayout(mockRegistryClient, mockLayoutPath)

		testError = errors.New("failed")
	})

	Describe("Find", func() {
		var (
			indexManifest *v1.IndexManifest
			im            image.Name
			dig           image.Digest
			err           error
			testHash      v1.Hash
		)

		BeforeEach(func() {
			indexManifest = &v1.IndexManifest{}
			mockImageIndex.IndexManifestReturns(indexManifest, nil)
			mockLayoutPath.ImageIndexReturns(mockImageIndex, nil)
			var err error
			im, err = image.NewName("testimage")
			Expect(err).NotTo(HaveOccurred())
			testHash, err = v1.NewHash("sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			dig, err = layout.Find(im)
		})

		Context("when the image is present", func() {
			BeforeEach(func() {
				indexManifest.Manifests = []v1.Descriptor{
					{
						Annotations: map[string]string{"org.opencontainers.image.ref.name": "testimage"},
						Digest:      testHash,
					},
				}
			})

			It("should find the image", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(dig.String()).To(Equal(testHash.String()))
			})
		})

		Context("when the image is absent", func() {
			It("should return a suitable error", func() {
				Expect(err).To(MatchError("image docker.io/library/testimage not found in layout"))
			})
		})

		Context("when accessing the image index returns an error", func() {
			BeforeEach(func() {
				mockLayoutPath.ImageIndexReturns(nil, testError)
			})

			It("should propagate the error", func() {
				Expect(err).To(MatchError(testError))
			})
		})

		Context("when accessing the index manifest returns an error", func() {
			BeforeEach(func() {
				mockImageIndex.IndexManifestReturns(nil, testError)
			})

			It("should propagate the error", func() {
				Expect(err).To(MatchError(testError))
			})
		})

		Context("when an image in the layout has an invalid name", func() {
			BeforeEach(func() {
				indexManifest.Manifests = []v1.Descriptor{
					{
						Annotations: map[string]string{"org.opencontainers.image.ref.name": ":"},
						Digest:      testHash,
					},
				}
			})

			It("should return a suitable error", func() {
				Expect(err).To(MatchError("invalid image reference: \":\""))
			})
		})
	})

	Describe("Push", func() {
		const testDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
		var (
			digest    image.Digest
			targetRef image.Name
			err       error
		)

		BeforeEach(func() {
			digest, err = image.NewDigest(testDigest)
			Expect(err).NotTo(HaveOccurred())
		})

		JustBeforeEach(func() {
			err = layout.Push(digest, targetRef)
		})

		Context("when the digest refers to a manifest", func() {
			var mockImage *ggcrfakes.FakeImage

			BeforeEach(func() {
				mockImage = &ggcrfakes.FakeImage{}
				mockLayoutPath.ImageIndexReturns(mockImageIndex, nil)
				mockImageIndex.ImageReturns(mockImage, nil)
			})

			It("should write the manifest", func() {
				h, er := v1.NewHash(testDigest)
				Expect(er).NotTo(HaveOccurred())
				Expect(mockImageIndex.ImageArgsForCall(0)).To(Equal(h))
				im, target := mockRegistryClient.WriteRemoteImageArgsForCall(0)
				Expect(im).To(Equal(mockImage))
				Expect(target).To(Equal(targetRef))
			})
		})

		Context("when the digest refers to an image index", func() {
			var mockImageIndex2 *ggcrfakes.FakeImageIndex

			BeforeEach(func() {
				mockImageIndex2 = &ggcrfakes.FakeImageIndex{}
				mockLayoutPath.ImageIndexReturns(mockImageIndex, nil)
				mockImageIndex.ImageReturns(nil, errors.New("some error"))
				mockImageIndex.ImageIndexReturns(mockImageIndex2, nil)
			})

			It("should write the manifest", func() {
				h, er := v1.NewHash(testDigest)
				Expect(er).NotTo(HaveOccurred())
				Expect(mockImageIndex.ImageArgsForCall(0)).To(Equal(h))
				Expect(mockImageIndex.ImageIndexArgsForCall(0)).To(Equal(h))
				idx, target := mockRegistryClient.WriteRemoteIndexArgsForCall(0)
				Expect(idx).To(Equal(mockImageIndex2))
				Expect(target).To(Equal(targetRef))
			})
		})

		Context("when the digest refers neither to a manifest nor an image index", func() {
			BeforeEach(func() {
				mockLayoutPath.ImageIndexReturns(mockImageIndex, nil)
				mockImageIndex.ImageReturns(nil, errors.New("image error"))
				mockImageIndex.ImageIndexReturns(nil, errors.New("index error"))
			})

			It("should return either the image lookup error or the index lookup error", func() {
				// Note: in practice the errors are identical, e.g. "could not find descriptor in index: sha256:..."
				Expect(err).To(Or(MatchError("image error"), MatchError("index error")))
			})
		})
	})
})
