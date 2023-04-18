package meta

// +gogo:genproto=true
// +gogo:deepcopy-gen=true
// +gogo:gengorm=true
// +gogo:gengorm:external=interfaces
// 资源元数据
type Resource struct {
	// +embedded
	// +protobuf.embed
	Meta `json:",inline" gorm:"embedded"`

	Spec string `json:"spec" gorm:"column:spec" protobuf:"bytes,1,opt,name=spec,proto3"`

	Label map[string]string `json:"label" protobuf:"bytes,2,rep,name=label,proto3"`

	//Labels dao.Array[string] `json:"labels" gorm:"column:labels;serializer:json"`
	//
	//Ann dao.Map[string, string] `json:"ann" gorm:"column:ann;serializer:json"`
	//
	//Subs dao.JSONArray[*Sub] `json:"subs" gorm:"column:subs;serializer:json"`
	//
	//SubMap dao.JSONMap[string, *Sub] `json:"subMap" gorm:"column:subMap;serializer:json"`
	//
	//Enable *bool `json:"enable" gorm:"column:enable"`
	//
	//Age int32 `json:"age" gorm:"column:age"`
}
