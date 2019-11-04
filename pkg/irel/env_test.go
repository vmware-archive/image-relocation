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

package irel

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CliVersion", func() {
	var version string

	JustBeforeEach(func() {
		version = CliVersion()
	})

	It("should default correctly", func() {
		Expect(version).To(Equal("unknown (unknown sha)"))
	})

	Context("when version is known", func() {
		BeforeEach(func() {
			cli_version = "0.0.0"
		})

		It("should include the version", func() {
			Expect(version).To(Equal("0.0.0 (unknown sha)"))
		})

		Context("when SHA is known", func() {
			BeforeEach(func() {
				cli_gitsha = "abc"
			})

			It("should include the SHA", func() {
				Expect(version).To(Equal("0.0.0 (abc)"))
			})

			Context("when repo is dirty", func() {
				BeforeEach(func() {
					cli_gitdirty = "dirty"
				})

				It("should note the repo is dirty", func() {
					Expect(version).To(Equal("0.0.0 (abc, with local modifications)"))
				})
			})
		})
	})
})
