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
	"fmt"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
	"github.com/spf13/cobra"
	"log"
)

func init() { Root.AddCommand(newCmdDigest()) }

func newCmdDigest() *cobra.Command {
	return &cobra.Command{
		Use:     "digest REF",
		Aliases: []string{"dig"},
		Short:   "Print content digest of an image",
		Args:    cobra.ExactArgs(1),
		Run:     digest,
	}
}

func digest(cmd *cobra.Command, args []string) {
	refStr := args[0]
	ref, err := image.NewName(refStr)
	if err != nil {
		log.Fatalf("invalid reference %q: %v", refStr, err)
	}

	regClient := ggcr.NewRegistryClient()
	dig, err := regClient.Digest(ref)
	if err != nil {
		log.Fatalf("digest failed: %v", err)
	}
	fmt.Printf("%s\n", dig)
}
