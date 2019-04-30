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

package registry

import (
	"errors"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/imagefakes"
)

var _ = Describe("Client", func() {
	var (
		cl              *client
		dig             image.Digest
		rawManifestSize int64
		err             error
		testError       error
		readArg         image.Name
		readResultImage v1.Image
		readResultErr   error
		writeArgImage   v1.Image
		writeArgName    image.Name
		writeResultErr  error
		hash            v1.Hash
		fakeImage       *imagefakes.FakeImage
	)

	BeforeEach(func() {
		cl = &client{}
		fakeImage = &imagefakes.FakeImage{}
		hash = createHash("sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
		fakeImage.DigestReturns(hash, nil)
		fakeImage.RawManifestReturns([]byte{0,1,2}, nil)
		readResultImage = fakeImage
		readResultErr = nil
		writeResultErr = nil
		cl.readRemoteImage = func(n image.Name) (v1.Image, error) {
			readArg = n
			return readResultImage, readResultErr
		}
		cl.writeRemoteImage = func(i v1.Image, n image.Name) error {
			writeArgImage = i
			writeArgName = n
			return writeResultErr
		}
		testError = errors.New("something bad happened")
	})

	Describe("Copy", func() {
		JustBeforeEach(func() {
			dig, rawManifestSize, err = cl.Copy(createName("source"), createName("target"))
		})

		Context("when no errors occur", func() {
			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
				_ = rawManifestSize
			})

			It("should copy the source repository to the target repository", func() {
				Expect(readArg.String()).To(Equal("docker.io/library/source"))
				Expect(writeArgImage).To(Equal(readResultImage))
				Expect(writeArgName.String()).To(Equal("docker.io/library/target"))
			})

			It("should return the correct digest", func() {
				Expect(dig.String()).To(Equal(hash.String()))
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
				fakeImage.DigestReturns(hash, testError)
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to read digest of image docker.io/library/source: something bad happened"))
			})
		})

		Context("when reading the raw manifest of the source image fails", func() {
			BeforeEach(func() {
				fakeImage.RawManifestReturns([]byte{}, testError)
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to get raw manifest of image docker.io/library/source: something bad happened"))
			})
		})

		Context("when writing the target image fails", func() {
			BeforeEach(func() {
				writeResultErr = testError
			})

			It("should return a corresponding error", func() {
				Expect(err).To(MatchError("failed to write image docker.io/library/target: something bad happened"))
			})
		})
	})

})

func createName(n string) image.Name {
	nm, err := image.NewName(n)
	Expect(err).NotTo(HaveOccurred())
	return nm
}

func createHash(h string) v1.Hash {
	hsh, err := v1.NewHash(h)
	Expect(err).NotTo(HaveOccurred())
	return hsh
}
