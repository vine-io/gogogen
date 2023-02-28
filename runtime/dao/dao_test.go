package dao

import (
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Kind string

type User struct {
}

func (u User) Scan(src any) error {
	//TODO implement me
	panic("implement me")
}

func (u User) Value() (driver.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (u User) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	//TODO implement me
	panic("implement me")
}

var _ JSONValue = (*User)(nil)

type Product struct {
	Labels Map[string, Kind]
	Items  Array[Kind]

	Users   JSONArray[*User]
	UserMap JSONMap[string, *User]
}
