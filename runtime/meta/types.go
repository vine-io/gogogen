// Copyright 2021 lack
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

// +gogo:genproto=true
// 资源元数据
type Meta struct {
	// 资源类型
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`
	// 资源版本
	APIVersion string `json:"apiVersion" protobuf:"bytes,2,opt,name=apiVersion"`
	// 资源名称
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// 资源的唯一ID
	UID string `json:"uid" protobuf:"bytes,4,opt,name=uid"`
	// 资源创建的时间戳
	CreationTimestamp int64 `json:"creationTimestamp" protobuf:"varint,5,opt,name=creationTimestamp"`
	// 资源更新的时间戳
	UpdateTimestamp int64 `json:"updateTimestamp" protobuf:"varint,6,opt,name=updateTimestamp"`
	// 资源删除的时间戳
	DeletionTimestamp int64 `json:"deletionTimestamp" protobuf:"varint,7,opt,name=deletionTimestamp"`
	// 资源标签
	Tags map[string]string `json:"tags" protobuf:"bytes,8,rep,name=tags"`
	// 资源注解
	Annotations map[string]string `json:"annotations" protobuf:"bytes,9,rep,name=annotations"`
}
