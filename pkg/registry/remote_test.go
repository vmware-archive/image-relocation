package registry

import (
	"errors"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/imagefakes"
)

var _ = Describe("remote utilities", func() {
	var (
		imageName image.Name
		mockImage *imagefakes.FakeImage
		testError error
		err       error
	)

	BeforeEach(func() {
		var err error
		imageName, err = image.NewName("imagename")
		Expect(err).NotTo(HaveOccurred())

		mockImage = &imagefakes.FakeImage{}
		h1, err := v1.NewHash("sha256:0000000000000000000000000000000000000000000000000000000000000000")
		Expect(err).NotTo(HaveOccurred())
		mockImage.DigestReturns(h1, nil)

		testError = errors.New("hard cheese")
	})

	Describe("readRemoteImage", func() {
		JustBeforeEach(func() {
			_, err = readRemoteImage(imageName)
		})

		BeforeEach(func() {
			var err error
			imageName, err = image.NewName("imagename")
			Expect(err).NotTo(HaveOccurred())

			// In most tests, keychain resolution succeeds
			resolveFunc = func(registry name.Registry) (authn.Authenticator, error) {
				return nil, nil
			}
		})

		Context("when keychain resolution fails", func() {
			BeforeEach(func() {
				resolveFunc = func(registry name.Registry) (authn.Authenticator, error) {
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
			err = writeRemoteImage(mockImage, imageName)
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
	})
})
