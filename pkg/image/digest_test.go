package image_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal/image-relocation/pkg/image"
)

var _ = Describe("Digest", func() {
	Describe("EmptyDigest", func() {
		It("should produce an empty string form", func() {
			Expect(image.EmptyDigest.String()).To(BeEmpty())
		})
	})

	Describe("NewDigest", func() {
		var (
			str    string
			digest image.Digest
		)

		JustBeforeEach(func() {
			digest = image.NewDigest(str)
		})

		Context("when the input string is empty", func() {
			BeforeEach(func() {
				str = ""
			})

			It("should produce an empty digest", func() {
				Expect(digest).To(Equal(image.EmptyDigest))
			})
		})

		Context("when the input string is non-empty", func() {
			BeforeEach(func() {
				str = "sha256:deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"
			})

			It("should produce a digest with the correct string form", func() {
				Expect(digest.String()).To(Equal(str))
			})
		})
	})
})
