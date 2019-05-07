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
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/spf13/cobra"
	"log"
)

func newCmdLayoutPush() *cobra.Command {
	return &cobra.Command{
		Use:     "push LAYOUT_PATH CONTENT_DIGEST REF",
		Short:   "copy an image with a given content digest from an OCI image layout to a remote repository",
		Args:    cobra.ExactArgs(3),
		Run:     layoutPush,
	}
}

func layoutPush(cmd *cobra.Command, args []string) {
	layoutPath, digStr, refStr := args[0], args[1], args[2]
	ref, err := image.NewName(refStr)
	if err != nil {
		log.Fatalf("invalid reference %q: %v", refStr, err)
	}

	dig, err := image.NewDigest(digStr)
	if err != nil {
		log.Fatalf("invalid digest %q: %v", digStr, err)
	}

	regClient := registry.NewRegistryClient()
	layout, err := regClient.ReadLayout(layoutPath)
	if err != nil {
		log.Fatalf("failed to access OCI image layout: %v", err)
	}

	err = layout.Push(dig, ref)
	if err != nil {
		log.Fatalf("push failed: %v", err)
	}
	fmt.Printf("wrote image with digest %s from OCI image layout at %s to %s\n", dig, layoutPath, ref)
}
