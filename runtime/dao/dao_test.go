package dao

import (
	"context"
	"database/sql/driver"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Kind string

type User[V _] struct {
}

func (u User[V]) Scan(src any) error {
	//TODO implement me
	panic("implement me")
}

func (u User[V]) Value() (driver.Value, error) {
	//TODO implement me
	panic("implement me")
}

func (u User[V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	//TODO implement me
	panic("implement me")
}

func GetById[R ~int](ctx context.Context, id R) {

}

var _ JSONValue = (*User)(nil)

type Product struct {
	Labels Map[string, Kind]
	Items  Array[Kind]

	Users   JSONArray[*User[int]]
	UserMap JSONMap[string, *User]
}
