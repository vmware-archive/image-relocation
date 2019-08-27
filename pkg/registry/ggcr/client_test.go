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
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/registryfakes"
)

var _ = Describe("Client", func() {
	var (
		cl              *client
		outputDig       image.Digest
		rawManifestSize int64
		err             error
		testError       error
		readArg         image.Name
		readResultImage registry.Image
		readResultErr   error
		writeArgName    image.Name
		dig             image.Digest
		dig2             image.Digest
		fakeImage       *registryfakes.FakeImage
	)

	BeforeEach(func() {
		cl = &client{}
		fakeImage = &registryfakes.FakeImage{}
		dig = createDigest("sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
		dig2 = createDigest("sha256:afabcafeafabcafeafabcafeafabcafeafabcafeafabcafeafabcafeafabcafe")
		fakeImage.DigestReturns(dig, nil)
		fakeImage.WriteStub = func(name image.Name) (digest image.Digest, i int64, e error) {
			writeArgName = name
			return dig, 3, nil
		}
		readResultImage = fakeImage
		readResultErr = nil
		cl.readRemoteImage = func(n image.Name) (registry.Image, error) {
			readArg = n
			return readResultImage, readResultErr
		}
		testError = errors.New("something bad happened")
	})

	Describe("Copy", func() {
		JustBeforeEach(func() {
			outputDig, rawManifestSize, err = cl.Copy(createName("source"), createName("target"))
		})

		Context("when no errors occur", func() {
			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
				_ = rawManifestSize
			})

			It("should copy the source repository to the target repository", func() {
				Expect(readArg.String()).To(Equal("docker.io/library/source"))
				Expect(writeArgName.String()).To(Equal("docker.io/library/target"))
			})

			It("should return the correct digest", func() {
				Expect(outputDig.String()).To(Equal(dig.String()))
			})

			It("should return the correct raw manifest size", func() {
				Expect(rawManifestSize).To(Equal(int64(3)))
			})
		})

		Context("when reading the source image fails", func() {
			BeforeEach(func() {
				readResultImage = nil
				readResultErr = testError
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to read image docker.io/library/source: something bad happened"))
			})
		})

		Context("when reading the source image digest fails", func() {
			BeforeEach(func() {
				fakeImage.DigestReturns(dig, testError)
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to read digest of image docker.io/library/source: something bad happened"))
			})
		})

		Context("when writing the target image image fails", func() {
			BeforeEach(func() {
				fakeImage.WriteStub = func(name image.Name) (digest image.Digest, i int64, e error) {
					return image.EmptyDigest, 0, testError
				}
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to write image docker.io/library/target: something bad happened"))
			})
		})

		Context("when writing the target image produces a distinct digest", func() {
			BeforeEach(func() {
				fakeImage.WriteStub = func(name image.Name) (digest image.Digest, i int64, e error) {
					return dig2, 0, nil
				}
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError(fmt.Sprintf("failed to preserve digest of image docker.io/library/source: source digest %v, target digest %v", dig, dig2)))
			})
		})
	})
})

func createName(n string) image.Name {
	nm, err := image.NewName(n)
	Expect(err).NotTo(HaveOccurred())
	return nm
}

func createDigest(h string) image.Digest {
	dig, err := image.NewDigest(h)
	Expect(err).NotTo(HaveOccurred())
	return dig
}
