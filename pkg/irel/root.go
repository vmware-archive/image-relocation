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
	"log"

	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
	"github.com/pivotal/image-relocation/pkg/transport"

	"github.com/spf13/cobra"
)

var (
	caCertPaths   []string
	skipTLSVerify bool

	// Root is the root of the tree of irel commands
	Root = &cobra.Command{
		Use:               "irel",
		Short:             "irel is a tool for relocating container images",
		Run:               func(cmd *cobra.Command, _ []string) { cmd.Usage() },
		DisableAutoGenTag: true,
	}
)

func init() {
	Root.PersistentFlags().StringSliceVarP(&caCertPaths, "ca-cert-path", "", nil, "Path to CA certificate for verifying registry TLS certificates (can be repeated for multiple certificates)")
	Root.PersistentFlags().BoolVarP(&skipTLSVerify, "skip-tls-verify", "", false, "Skip TLS certificate verification for registries")
}

func mustGetRegistryClient() registry.Client {
	tport, err := transport.NewHttpTransport(caCertPaths, skipTLSVerify)
	if err != nil {
		log.Fatal(err)
	}

	return ggcr.NewRegistryClient(ggcr.WithTransport(tport))
}
