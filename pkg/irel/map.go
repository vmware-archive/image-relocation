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
	"github.com/pivotal/image-relocation/pkg/pathmapping"
	"github.com/spf13/cobra"
)

func init() { Root.AddCommand(newCmdMap()) }

func newCmdMap() *cobra.Command {
	var repoPrefix string
	cmd := &cobra.Command{
		Use:   "map REF",
		Short: "Map an image reference to a relocated reference",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			pathMapping(repoPrefix, args)
		},
	}
	cmd.Flags().StringVarP(&repoPrefix, "repository-prefix", "r", "", "base value to which an image name is appended to create the full repository")
	cmd.MarkFlagRequired("repository-prefix")

	return cmd
}

func pathMapping(repoPrefix string, args []string) {
	refStr := args[0]
	ref, err := image.NewName(refStr)
	if err != nil {
		log.Fatalf("invalid reference %q: %v", refStr, err)
	}

	mapped, err := pathmapping.FlattenRepoPathPreserveTagDigest(repoPrefix, ref)
	if err != nil {
		log.Fatalf("path flattening failed: %v", err)
	}
	fmt.Printf("%s\n", mapped)
}
