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

package transport_test

import (
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pivotal/image-relocation/pkg/transport"
)

var _ = Describe("NewHttpTransport", func() {
	var (
		certs              []string
		insecureSkipVerify bool
		httpTransport      *http.Transport
		err                error
	)

	BeforeEach(func() {
		certs = []string{}
		insecureSkipVerify = false
	})

	JustBeforeEach(func() {
		httpTransport, err = transport.NewHttpTransport(certs, insecureSkipVerify)
	})

	Context("when skipping TLS certificate verification is set to false", func() {
		BeforeEach(func() {
			insecureSkipVerify = false
		})

		It("should not skip TLS verification", func() {
			Expect(err).NotTo(HaveOccurred())
			if httpTransport.TLSClientConfig != nil {
				Expect(httpTransport.TLSClientConfig.InsecureSkipVerify).To(BeFalse())
			}
		})

		Context("when CA certs are provided", func() {
			BeforeEach(func() {
				certs = []string{"testdata/ca.crt"}
			})

			It("should not skip TLS verification", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(httpTransport.TLSClientConfig.InsecureSkipVerify).To(BeFalse())
			})

			It("should use the provided CA certs", func() {
				// Check there is a subject with organization ACME
				Expect(findSubject(httpTransport.TLSClientConfig.RootCAs.Subjects(), "ACME")).To(BeTrue())
			})
		})

		Context("when an empty CA cert is provided", func() {
			BeforeEach(func() {
				certs = []string{"testdata/empty.crt"}
			})

			It("should return a suitable error", func() {
				Expect(err).To(MatchError(`could not append "testdata/empty.crt" to certificate pool`))
			})
		})

		Context("when a non-existent CA cert is provided", func() {
			BeforeEach(func() {
				certs = []string{"testdata/nosuch.crt"}
			})

			It("should return a suitable error", func() {
				Expect(err).To(MatchError(HavePrefix(`could not read certificates from "testdata/nosuch.crt":`)))
			})
		})
	})

	Context("when skipping TLS certificate verification is set to true", func() {
		BeforeEach(func() {
			insecureSkipVerify = true
		})

		It("should skip TLS verification", func() {
			Expect(err).NotTo(HaveOccurred())
			Expect(httpTransport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
		})

		Context("when CA certs are also provided", func() {
			BeforeEach(func() {
				certs = []string{"testdata/ca.crt"}
			})

			It("should still skip TLS verification", func() {
				Expect(err).NotTo(HaveOccurred())
				Expect(httpTransport.TLSClientConfig.InsecureSkipVerify).To(BeTrue())
			})
		})
	})
})

func findSubject(subjects [][]byte, org string) bool {
	found := false
	for _, subject := range subjects {
		if strings.Contains(string(subject), org) {
			found = true
			break
		}
	}
	return found
}
