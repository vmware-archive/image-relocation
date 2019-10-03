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

// ClusterImageMapSpec defines the desired state of ClusterImageMap
type ClusterImageMapSpec struct {
	Map map[string]string `json:"map,omitempty"`
}

func (in *ClusterImageMapSpec) DeepCopyInto(out *ClusterImageMapSpec) {
	out.Map = make(map[string]string, len(in.Map))
	for k, v := range in.Map {
		out.Map[k] = v
	}
}

// ClusterImageMapStatus defines the observed state of ClusterImageMap
type ClusterImageMapStatus struct {
	Map map[string]string `json:"map,omitempty"`
}

func (in *ClusterImageMapStatus) DeepCopyInto(out *ClusterImageMapStatus) {
	out.Map = make(map[string]string, len(in.Map))
	for k, v := range in.Map {
		out.Map[k] = v
	}
}

// ClusterImageMap is the Schema for the imagemaps API
type ClusterImageMap struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterImageMapSpec   `json:"spec,omitempty"`
	Status ClusterImageMapStatus `json:"status,omitempty"`
}

// ClusterImageMapList contains a list of ClusterImageMap
type ClusterImageMapList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterImageMap `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterImageMap{}, &ClusterImageMapList{})
}
