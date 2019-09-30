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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields must have json tags for the fields to be serialized.

// ImageMapSpec defines the desired state of ImageMap
type ImageMapSpec struct {
	Map map[string]string `json:"map,omitempty"`
}

func (in *ImageMapSpec) DeepCopyInto(out *ImageMapSpec) {
	out.Map = make(map[string]string, len(in.Map))
	for k, v := range in.Map {
		out.Map[k] = v
	}
}

// ImageMapStatus defines the observed state of ImageMap
type ImageMapStatus struct {
	Map map[string]string `json:"map,omitempty"`
}

func (in *ImageMapStatus) DeepCopyInto(out *ImageMapStatus) {
	out.Map = make(map[string]string, len(in.Map))
	for k, v := range in.Map {
		out.Map[k] = v
	}
}

// ImageMap is the Schema for the imagemaps API
type ImageMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageMapSpec   `json:"spec,omitempty"`
	Status ImageMapStatus `json:"status,omitempty"`
}

// ImageMapList contains a list of ImageMap
type ImageMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ImageMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ImageMap{}, &ImageMapList{})
}
