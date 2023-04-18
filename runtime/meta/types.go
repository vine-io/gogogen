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

// +gogo:deepcopy-gen=package
package meta

// +gogo:deepcopy=true
// +gogo:deepcopy:interfaces=github.com/vine-io/apimachinery/runtime.Object
// +gogo:genproto=true
// +gogo:gengorm=true
// +gogo:gengorm:external=false
// 资源元数据
type Meta struct {
	// 资源类型
	Kind string `json:"kind" gorm:"column:kind" protobuf:"bytes,1,opt,name=kind,proto3"`
	// 资源版本
	APIVersion string `json:"apiVersion" gorm:"column:apiVersion" protobuf:"bytes,2,opt,name=apiVersion,proto3"`
	// 资源名称
	Name string `json:"name" gorm:"column:name" protobuf:"bytes,3,opt,name=name,proto3"`
	// 资源的唯一ID
	// +primaryKey
	UID string `json:"uid" gorm:"column:uid;primaryKey" protobuf:"bytes,4,opt,name=uid,proto3"`
	// 资源创建的时间戳
	CreationTimestamp int64 `json:"creationTimestamp" gorm:"column:creationTimestamp" protobuf:"varint,5,opt,name=creationTimestamp,proto3"`
	// 资源更新的时间戳
	UpdateTimestamp int64 `json:"updateTimestamp" gorm:"column:updateTimestamp" protobuf:"varint,6,opt,name=updateTimestamp,proto3"`
	// 资源删除的时间戳
	DeletionTimestamp int64 `json:"deletionTimestamp" gorm:"column:deletionTimestamp" protobuf:"varint,7,opt,name=deletionTimestamp,proto3"`
	// 资源标签
	Tags map[string]string `json:"tags" gorm:"column:tags;serializer:json" protobuf:"bytes,8,rep,name=tags,proto3"`
	// 资源注解
	Annotations map[string]string `json:"annotations" gorm:"column:annotations;serializer:json" protobuf:"bytes,9,rep,name=annotations,proto3"`
}

// +gogo:deepcopy-gen=true
// +gogo:genproto=true
// +gogo:gengorm=true
// +gogo:gengorm:external=false
type Sub struct {
	Name string `json:"name" gorm:"column:name" protobuf:"bytes,1,opt,name=name,proto3"`
	// +primaryKey
	Age int32 `json:"age" gorm:"column:age;primaryKey" protobuf:"varint,2,opt,name=age,proto3"`
}
