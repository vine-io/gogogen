// Copyright 2020 The vine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

type Meta struct {
	// 资源的唯一ID
	UID string `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	// 资源创建的时间戳
	CreationTimestamp int64 `json:"creationTimestamp" protobuf:"varint,2,opt,name=creationTimestamp"`
	// 资源更新的时间戳
	DeletionTimestamp int64 `json:"deletionTimestamp" protobuf:"varint,3,opt,name=deletionTimestamp"`
}