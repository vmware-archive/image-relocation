package ggcr

import (
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcrfakes"
)

var _ = Describe("remote utilities", func() {
	var (
		imageName  image.Name
		mockImage  *ggcrfakes.FakeImage
		testDigest string
		testError  error
		err        error
	)

	BeforeEach(func() {
		var err error
		imageName, err = image.NewName("imagename")
		Expect(err).NotTo(HaveOccurred())

		testDigest = "sha256:0000000000000000000000000000000000000000000000000000000000000000"
		mockImage = &ggcrfakes.FakeImage{}
		h1, err := v1.NewHash(testDigest)
		Expect(err).NotTo(HaveOccurred())
		mockImage.DigestReturns(h1, nil)

		testError = errors.New("hard cheese")
	})

	// FIXME: get coverage back up

	Describe("readRemoteImage", func() {
		JustBeforeEach(func() {
			_, err = readRemoteImage(nil, nil, nil)(imageName)
		})

		BeforeEach(func() {
			var err error
			imageName, err = image.NewName("imagename")
			Expect(err).NotTo(HaveOccurred())

			// In most tests, keychain resolution succeeds
			resolveFunc = func(authn.Resource) (authn.Authenticator, error) {
				return nil, nil
			}
		})

		Context("when keychain resolution fails", func() {
			BeforeEach(func() {
				resolveFunc = func(authn.Resource) (authn.Authenticator, error) {
					return nil, testError
				}
			})

			It("should return the error", func() {
				Expect(err).To(Equal(testError))
			})
		})

		Context("when the image name is empty", func() {
			BeforeEach(func() {
				imageName = image.EmptyName
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("empty image name invalid"))
			})
		})
	})

	Describe("writeRemoteImage", func() {
		JustBeforeEach(func() {
			err = writeRemoteImage(nil)(mockImage, imageName)
		})

		Context("when writing to the repository succeeds", func() {
			BeforeEach(func() {
				repoWriteFunc = func(ref name.Reference, img v1.Image, options ...remote.Option) error {
					return nil
				}
			})

			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when writing to the repository return an error", func() {
			BeforeEach(func() {
				repoWriteFunc = func(ref name.Reference, img v1.Image, options ...remote.Option) error {
					return testError
				}
			})

			It("should return the error", func() {
				Expect(err).To(Equal(testError))
			})
		})

		Context("when the image name is empty", func() {
			BeforeEach(func() {
				imageName = image.EmptyName
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("empty image name invalid"))
			})
		})

		Context("when the image name is both tagged and digested", func() {
			var writeRef name.Reference
			BeforeEach(func() {
				imageName, err = image.NewName(fmt.Sprintf("example.com/eg:1@%s", testDigest))
				Expect(err).NotTo(HaveOccurred())
				repoWriteFunc = func(ref name.Reference, img v1.Image, options ...remote.Option) error {
					writeRef = ref
					return nil
				}
			})

			It("should discard the digest from the written reference", func() {
				Expect(writeRef.String()).To(Equal("example.com/eg:1"))
			})
		})
	})

	Describe("writeRemoteIndex", func() {
		var mockIndex *ggcrfakes.FakeImageIndex

		BeforeEach(func() {
			mockIndex = &ggcrfakes.FakeImageIndex{}
		})

		JustBeforeEach(func() {
			err = writeRemoteIndex(nil)(mockIndex, imageName)
		})

		Context("when writing to the repository succeeds", func() {
			BeforeEach(func() {
				repoIndexWriteFunc = func(ref name.Reference, ii v1.ImageIndex, options ...remote.Option) error {
					return nil
				}
			})

			It("should succeed", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when writing to the repository return an error", func() {
			BeforeEach(func() {
				repoIndexWriteFunc = func(ref name.Reference, ii v1.ImageIndex, options ...remote.Option) error {
					return testError
				}
			})

			It("should return the error", func() {
				Expect(err).To(Equal(testError))
			})
		})

		Context("when the image name is empty", func() {
			BeforeEach(func() {
				imageName = image.EmptyName
			})

			It("should return an error", func() {
				Expect(err).To(MatchError("empty image name invalid"))
			})
		})

		Context("when the image name is both tagged and digested", func() {
			var writeRef name.Reference
			BeforeEach(func() {
				imageName, err = image.NewName(fmt.Sprintf("example.com/eg:1@%s", testDigest))
				Expect(err).NotTo(HaveOccurred())
				repoIndexWriteFunc = func(ref name.Reference, ii v1.ImageIndex, options ...remote.Option) error {
					writeRef = ref
					return nil
				}
			})

			It("should discard the digest from the written reference", func() {
				Expect(writeRef.String()).To(Equal("example.com/eg:1"))
			})
		})
	})
})
