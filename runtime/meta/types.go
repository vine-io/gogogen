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

import "github.com/vine-io/gogogen/runtime/dao"

// +gogo:deepcopy-gen=true
// +gogo:genproto=true
// +gogo:gengorm=true
// 资源元数据
type Meta struct {
	// 资源类型
	Kind string `json:"kind" gorm:"column:kind" protobuf:"bytes,1,opt,name=kind"`
	// 资源版本
	APIVersion string `json:"apiVersion" gorm:"column:apiVersion" protobuf:"bytes,10,opt,name=apiVersion"`
	// 资源名称
	Name string `json:"name" gorm:"column:name" protobuf:"bytes,3,opt,name=name"`
	// 资源的唯一ID
	// +primaryKey
	UID string `json:"uid" gorm:"column:uid;primaryKey" protobuf:"bytes,11,opt,name=uid"`
	// 资源创建的时间戳
	CreationTimestamp int64 `json:"creationTimestamp" gorm:"column:creationTimestamp" protobuf:"varint,5,opt,name=creationTimestamp"`
	// 资源更新的时间戳
	UpdateTimestamp int64 `json:"updateTimestamp" gorm:"column:updateTimestamp" protobuf:"varint,6,opt,name=updateTimestamp"`
	// 资源删除的时间戳
	DeletionTimestamp int64 `json:"deletionTimestamp" gorm:"column:deletionTimestamp" protobuf:"varint,7,opt,name=deletionTimestamp"`
	// 资源标签
	Tags map[string]string `json:"tags" gorm:"column:tags;serializer:json" protobuf:"bytes,8,rep,name=tags"`
	// 资源注解
	Annotations map[string]string `json:"annotations" gorm:"column:annotations;serializer:json" protobuf:"bytes,9,rep,name=annotations"`
}

// +gogo:deepcopy-gen=true
// +gogo:genproto=true
// +gogo:gengorm=true
type Sub struct {
	Name string `json:"name" gorm:"column:name" protobuf:"bytes,1,opt,name=name"`
	// +primaryKey
	Age int32 `json:"age" gorm:"column:age;primaryKey" protobuf:"varint,2,opt,name=age"`
}

// +gogo:genproto=true
// +gogo:deepcopy-gen=true
// +gogo:gengorm=true
// 资源元数据
type Resource struct {
	// +embedded
	// +protobuf.embed
	Meta `json:",inline" gorm:"embedded" protobuf:"bytes,1,req,name=meta"`

	Spec string `json:"spec" gorm:"column:spec" protobuf:"bytes,2,opt,name=spec"`

	Labels dao.Array[string] `json:"labels" gorm:"column:labels;serializer:json" protobuf:"bytes,3,rep,name=labels"`

	Ann dao.Map[string, string] `json:"ann" gorm:"column:ann;serializer:json" protobuf:"bytes,4,rep,name=ann"`

	Subs dao.JSONArray[*Sub] `json:"subs" gorm:"column:subs;serializer:json" protobuf:"bytes,5,rep,name=subs"`

	SubMap dao.JSONMap[string, *Sub] `json:"subMap" gorm:"column:subMap;serializer:json" protobuf:"bytes,6,rep,name=subMap"`

	Enable *bool `json:"enable" gorm:"column:enable" protobuf:"varint,7,opt,name=enable"`

	Age int32 `json:"age" gorm:"column:age" protobuf:"varint,8,opt,name=age"`
}
