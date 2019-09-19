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
	"log"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
	"github.com/spf13/cobra"
)

func newCmdLayoutAdd() *cobra.Command {
	return &cobra.Command{
		Use:   "add LAYOUT_PATH REF",
		Short: "copy an image from a remote repository to an OCI image layout",
		Args:  cobra.ExactArgs(2),
		Run:   layoutAdd,
	}
}

func layoutAdd(cmd *cobra.Command, args []string) {
	layoutPath, refStr := args[0], args[1]
	ref, err := image.NewName(refStr)
	if err != nil {
		log.Fatalf("invalid reference %q: %v", refStr, err)
	}

	regClient := ggcr.NewRegistryClient()
	layout, err := regClient.ReadLayout(layoutPath)
	if err != nil {
		layout, err = regClient.NewLayout(layoutPath)
		if err != nil {
			log.Fatalf("failed to create OCI image layout: %v", err)
		}
	}

	dig, err := layout.Add(ref)
	if err != nil {
		log.Fatalf("add failed: %v", err)
	}
	fmt.Printf("added %s with digest %s to OCI image layout at %s\n", ref, dig, layoutPath)
}
