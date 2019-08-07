package registry

import (
	"errors"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
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
		var (
			mockImage2 *imagefakes.FakeImage

			img v1.Image
		)

		JustBeforeEach(func() {
			img, err = readRemoteImage(imageName)
		})

		BeforeEach(func() {
			var err error
			imageName, err = image.NewName("imagename")
			Expect(err).NotTo(HaveOccurred())

			mockImage2 = &imagefakes.FakeImage{}
			h2, err := v1.NewHash("sha256:1111111111111111111111111111111111111111111111111111111111111111")
			Expect(err).NotTo(HaveOccurred())
			mockImage2.DigestReturns(h2, nil)

			Expect(mockImage).NotTo(Equal(mockImage2)) // crucial for the correctness of certain tests

			// In most tests, keychain resolution succeeds
			resolveFunc = func(registry name.Registry) (authn.Authenticator, error) {
				return nil, nil
			}
		})

		Context("when the daemon returns an image", func() {
			BeforeEach(func() {
				daemonImageFunc = func(ref name.Reference, options ...daemon.ImageOption) (v1.Image, error) {
					return mockImage, nil
				}
			})

			It("should return the image", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(img).To(Equal(mockImage))
			})

			Context("when the remote repository returns an image", func() {
				BeforeEach(func() {
					repoImageFunc = func(ref name.Reference, options ...remote.Option) (v1.Image, error) {
						return mockImage2, nil
					}
				})

				It("should return the image from the daemon", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(img).To(Equal(mockImage))
				})
			})

			Context("when the remote repository returns an error", func() {
				BeforeEach(func() {
					repoImageFunc = func(ref name.Reference, options ...remote.Option) (v1.Image, error) {
						return nil, testError
					}
				})

				It("should return the image from the daemon", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(img).To(Equal(mockImage))
				})
			})

		})

		Context("when the daemon returns an error", func() {
			BeforeEach(func() {
				daemonImageFunc = func(ref name.Reference, options ...daemon.ImageOption) (v1.Image, error) {
					return nil, testError
				}
			})

			Context("when the remote repository returns an image", func() {
				BeforeEach(func() {
					repoImageFunc = func(ref name.Reference, options ...remote.Option) (v1.Image, error) {
						return mockImage, nil
					}
				})

				It("should return the image", func() {
					Expect(err).NotTo(HaveOccurred())
					Expect(img).To(Equal(mockImage))
				})
			})

			Context("when the remote repository returns an error", func() {
				BeforeEach(func() {
					repoImageFunc = func(ref name.Reference, options ...remote.Option) (v1.Image, error) {
						return nil, errors.New("repo error")
					}
				})

				It("should return a combined error", func() {
					Expect(err).To(MatchError("reading remote image docker.io/library/imagename:latest failed: repo error; attempting to read from daemon also failed: hard cheese"))
				})
			})
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
