package dao

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type JSONValue interface {
	sql.Scanner
	driver.Valuer
	GormDBDataType(db *gorm.DB, field *schema.Field) string
}

type Builtin interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~complex64 | ~complex128 |
		~float32 | ~float64 |
		~string | ~byte | ~rune | ~uintptr
}

type Array[V Builtin] []V

// Value return json value, implement driver.Valuer interface
func (m *Array[V]) Value() (driver.Value, error) {
	return GetValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Array[V]) Scan(value any) error {
	return ScanValue(value, m)
}

func (m *Array[V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return GetGormDBDataType(db, field)
}

type JSONArray[V JSONValue] []V

// Value return json value, implement driver.Valuer interface
func (m *JSONArray[V]) Value() (driver.Value, error) {
	return GetValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONArray[V]) Scan(value any) error {
	return ScanValue(value, m)
}

func (m *JSONArray[V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return GetGormDBDataType(db, field)
}

type Map[K comparable, V Builtin] map[K]V

// Value return json value, implement driver.Valuer interface
func (m *Map[K, V]) Value() (driver.Value, error) {
	return GetValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Map[K, V]) Scan(value any) error {
	return ScanValue(value, m)
}

func (m *Map[K, V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return GetGormDBDataType(db, field)
}

type JSONMap[K comparable, V JSONValue] map[K]V

// Value return json value, implement driver.Valuer interface
func (m *JSONMap[K, V]) Value() (driver.Value, error) {
	return GetValue(m)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONMap[K, V]) Scan(value any) error {
	return ScanValue(value, m)
}

func (m *JSONMap[K, V]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return GetGormDBDataType(db, field)
}

// GetValue return json value, implement driver.Valuer interface
func GetValue(m any) (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	b, err := json.Marshal(m)
	return string(b), err
}

// ScanValue scan value into Jsonb, implements sql.Scanner interface
func ScanValue(value, m any) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	return json.Unmarshal(bytes, &m)
}

func GetGormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "mysql", "sqlite", "splite3":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
