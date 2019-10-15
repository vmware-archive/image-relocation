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

package ggcr

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr/path/pathfakes"
	"github.com/pivotal/image-relocation/pkg/registry/ggcrfakes"
)

var _ = Describe("Image", func() {
	var (
		nm         image.Name
		testDigest image.Digest
		testHash   v1.Hash
		testErr    error
		err        error
	)

	BeforeEach(func() {
		nm, err = image.NewName("ubuntu")
		Expect(err).NotTo(HaveOccurred())

		const sha = "sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
		testDigest, err = image.NewDigest(sha)
		Expect(err).NotTo(HaveOccurred())
		testHash, err = v1.NewHash(sha)
		Expect(err).NotTo(HaveOccurred())

		testErr = errors.New("wat")
	})

	Describe("imageManifest", func() {
		var (
			im                 *imageManifest
			mockManifest       *ggcrfakes.FakeImage
			manifestWriterStub manifestWriter
		)

		BeforeEach(func() {
			mockManifest = &ggcrfakes.FakeImage{}
		})

		JustBeforeEach(func() {
			im = newImageFromManifest(mockManifest, manifestWriterStub)
		})

		Describe("Digest", func() {
			Context("when the digest is available", func() {
				BeforeEach(func() {
					mockManifest.DigestReturns(testHash, nil)
				})

				It("should return the digest", func() {
					d, err := im.Digest()
					Expect(err).NotTo(HaveOccurred())
					Expect(d).To(Equal(testDigest))
				})
			})

			Context("when the digest is not available", func() {
				BeforeEach(func() {
					mockManifest.DigestReturns(testHash, testErr)
				})

				It("should return a suitable error", func() {
					_, err := im.Digest()
					Expect(err).To(MatchError(testErr))
				})
			})
		})

		Describe("Write", func() {
			var (
				dig          image.Digest
				sz           int64
				writtenImage v1.Image
				writtenName  image.Name
				writeErr     error
			)

			BeforeEach(func() {
				writeErr = nil
				manifestWriterStub = func(i v1.Image, n image.Name) error {
					writtenImage = i
					writtenName = n
					return writeErr
				}
			})

			JustBeforeEach(func() {
				dig, sz, err = im.Write(nm)
			})

			Context("when the digest is available", func() {
				var testRawManifest []byte

				BeforeEach(func() {
					mockManifest.DigestReturns(testHash, nil)
					testRawManifest = []byte{0, 1}
					mockManifest.RawManifestReturns(testRawManifest, nil)
				})

				It("should return the digest", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(dig).To(Equal(testDigest))
					Expect(sz).To(Equal(int64(len(testRawManifest))))
				})

				It("should write the image", func() {
					Expect(writtenImage).To(Equal(mockManifest))
					Expect(writtenName).To(Equal(nm))
				})

				Context("when writing fails", func() {
					BeforeEach(func() {
						writeErr = testErr
					})

					It("should return a suitable error", func() {
						Expect(err).To(MatchError("failed to write image docker.io/library/ubuntu: wat"))
					})
				})

				Context("when the raw manifest is not available", func() {
					BeforeEach(func() {
						mockManifest.RawManifestReturns(nil, testErr)
					})

					It("should return a suitable error", func() {
						Expect(err).To(MatchError("failed to get raw manifest of image: wat"))
					})
				})
			})

			Context("when the digest is not available", func() {
				BeforeEach(func() {
					mockManifest.DigestReturns(testHash, testErr)
				})

				It("should return a suitable error", func() {
					Expect(err).To(MatchError("failed to read digest of image: wat"))
				})
			})
		})

		Describe("appendToLayout", func() {
			var mockLayoutPath *pathfakes.FakeLayoutPath

			BeforeEach(func() {
				mockLayoutPath = &pathfakes.FakeLayoutPath{}
			})

			JustBeforeEach(func() {
				err = im.appendToLayout(mockLayoutPath)
			})

			Context("when appending the image succeeds", func() {
				It("should pass the correct image", func() {
					Expect(err).NotTo(HaveOccurred())
					mfst, opts := mockLayoutPath.AppendImageArgsForCall(0)
					Expect(mfst).To(Equal(mockManifest))
					Expect(opts).To(BeEmpty())
				})
			})

			Context("when appending the image fails", func() {
			    BeforeEach(func() {
			        mockLayoutPath.AppendImageReturns(testErr)
			    })

			    It("should return the error", func() {
			        Expect(err).To(MatchError(testErr))
			    })
			})
		})
	})

	Describe("imageIndex", func() {
		var (
			im              *imageIndex
			mockIndex       *ggcrfakes.FakeImageIndex
			indexWriterStub indexWriter
		)

		BeforeEach(func() {
			mockIndex = &ggcrfakes.FakeImageIndex{}
		})

		JustBeforeEach(func() {
			im = newImageFromIndex(mockIndex, indexWriterStub)
		})

		Describe("Digest", func() {
			Context("when the digest is available", func() {
				BeforeEach(func() {
					mockIndex.DigestReturns(testHash, nil)
				})

				It("should return the digest", func() {
					d, err := im.Digest()
					Expect(err).NotTo(HaveOccurred())
					Expect(d).To(Equal(testDigest))
				})
			})

			Context("when the digest is not available", func() {
				BeforeEach(func() {
					mockIndex.DigestReturns(testHash, testErr)
				})

				It("should return the error", func() {
					_, err := im.Digest()
					Expect(err).To(MatchError(testErr))
				})
			})
		})

		Describe("Write", func() {
			var (
				dig          image.Digest
				sz           int64
				writtenImage v1.ImageIndex
				writtenName  image.Name
				writeErr     error
			)

			BeforeEach(func() {
				writeErr = nil
				indexWriterStub = func(i v1.ImageIndex, n image.Name) error {
					writtenImage = i
					writtenName = n
					return writeErr
				}
			})

			JustBeforeEach(func() {
				dig, sz, err = im.Write(nm)
			})

			Context("when the digest is available", func() {
				var testRawManifest []byte

				BeforeEach(func() {
					mockIndex.DigestReturns(testHash, nil)
					testRawManifest = []byte{0, 1}
					mockIndex.RawManifestReturns(testRawManifest, nil)
				})

				It("should return the digest", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(dig).To(Equal(testDigest))
					Expect(sz).To(Equal(int64(len(testRawManifest))))
				})

				It("should write the image", func() {
					Expect(writtenImage).To(Equal(mockIndex))
					Expect(writtenName).To(Equal(nm))
				})

				Context("when writing fails", func() {
					BeforeEach(func() {
						writeErr = testErr
					})

					It("should return a suitable error", func() {
						Expect(err).To(MatchError("failed to write image index docker.io/library/ubuntu: wat"))
					})
				})

				Context("when the raw manifest is not available", func() {
					BeforeEach(func() {
						mockIndex.RawManifestReturns(nil, testErr)
					})

					It("should return a suitable error", func() {
						Expect(err).To(MatchError("failed to get raw manifest of image index: wat"))
					})
				})
			})

			Context("when the digest is not available", func() {
				BeforeEach(func() {
					mockIndex.DigestReturns(testHash, testErr)
				})

				It("should return a suitable error", func() {
					Expect(err).To(MatchError("failed to read digest of image index: wat"))
				})
			})
		})

		Describe("appendToLayout", func() {
			var mockLayoutPath *pathfakes.FakeLayoutPath

			BeforeEach(func() {
				mockLayoutPath = &pathfakes.FakeLayoutPath{}
			})

			JustBeforeEach(func() {
				err = im.appendToLayout(mockLayoutPath)
			})

			Context("when appending the image succeeds", func() {
				It("should pass the correct image", func() {
					Expect(err).NotTo(HaveOccurred())
					mfst, opts := mockLayoutPath.AppendIndexArgsForCall(0)
					Expect(mfst).To(Equal(mockIndex))
					Expect(opts).To(BeEmpty())
				})
			})

			Context("when appending the image fails", func() {
				BeforeEach(func() {
					mockLayoutPath.AppendIndexReturns(testErr)
				})

				It("should return the error", func() {
					Expect(err).To(MatchError(testErr))
				})
			})
		})
	})
})
