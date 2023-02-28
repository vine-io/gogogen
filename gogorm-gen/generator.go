// Copyright 2020 lack
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

package goproto_gen

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/vine-io/gogogen/gogenerator/generator"
	"github.com/vine-io/gogogen/gogenerator/namer"
	"github.com/vine-io/gogogen/gogenerator/types"
	"github.com/vine-io/gogogen/util/log"
)

// genGormIDL produces a gorm IDL.
type genGormIDL struct {
	generator.DefaultGen
	localPackage   types.Name
	localGoPackage types.Name
	imports        namer.ImportTracker

	generateAll    bool
	omitFieldTypes map[types.Name]struct{}
}

func (g *genGormIDL) PackageVars(c *generator.Context) []string {
	return []string{
		fmt.Sprintf("package %s", g.localGoPackage.Name),
	}
}
func (g *genGormIDL) Filename() string { return g.OptionalName + ".go" }
func (g *genGormIDL) FileType() string { return "gormidl" }
func (g *genGormIDL) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		// The local namer returns the correct protobuf name for a proto type
		// in the context of a package
		"local": localNamer{g.localPackage},
	}
}

// Filter ignores types that are identified as not exportable.
func (g *genGormIDL) Filter(c *generator.Context, t *types.Type) bool {
	tagVals := types.ExtractCommentTags("+", t.CommentLines)["gorm"]
	if tagVals != nil {
		if tagVals[0] == "false" {
			// Type specified "false".
			return false
		}
		if tagVals[0] == "true" {
			// Type specified "true".
			return true
		}
		log.Fatalf(`Comment tag "gorm" must be true or false, found: %q`, tagVals[0])
	}
	if !g.generateAll {
		// We're not generating everything.
		return false
	}
	seen := map[*types.Type]bool{}
	ok := isGormable(seen, t)
	return ok
}

func isGormable(seen map[*types.Type]bool, t *types.Type) bool {
	if seen[t] {
		// be optimistic in the case of type cycles.
		return true
	}
	seen[t] = true
	switch t.Kind {
	case types.Builtin:
		return true
	case types.Alias:
		return isGormable(seen, t.Underlying)
	case types.Slice, types.Pointer:
		return isGormable(seen, t.Elem)
	case types.Map:
		return isGormable(seen, t.Key) && isGormable(seen, t.Elem)
	case types.Struct:
		if len(t.Members) == 0 {
			return true
		}
		for _, m := range t.Members {
			if isGormable(seen, m.Type) {
				return true
			}
		}
		return false
	case types.Func, types.Chan:
		return false
	case types.DeclarationOf, types.Unknown, types.Unsupported:
		return false
	case types.Interface:
		return false
	default:
		log.Warnf("WARNING: type %q is not portable: %s", t.Kind, t.Name)
		return false
	}
}

// isOptionalAlias should return true if the specified type has an underlying type
// (is an alias) of a map or slice and has the comment tag protobuf.nullable=true,
// indicating that the type should be nullable in protobuf.
func isOptionalAlias(t *types.Type) bool {
	if t.Underlying == nil || (t.Underlying.Kind != types.Map && t.Underlying.Kind != types.Slice) {
		return false
	}
	if extractBoolTagOrDie("gorm.nullable", t.CommentLines) == false {
		return false
	}
	return true
}

func (g *genGormIDL) Imports(c *generator.Context) (imports []string) {
	lines := []string{}
	// TODO: this could be expressed more cleanly
	for _, line := range g.imports.ImportLines() {
		lines = append(lines, line)
	}
	return lines
}

// GenerateType makes the body of a file implementing a set for type t.
func (g *genGormIDL) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	b := bodyGen{
		locator: &gormLocator{
			namer:    c.Namers["gorm"].(GormFromGoNamer),
			tracker:  g.imports,
			universe: c.Universe,

			localGoPackage: g.localGoPackage.Package,
		},
		localPackage: g.localPackage,

		omitFieldTypes: g.omitFieldTypes,

		t: t,
	}
	switch t.Kind {
	case types.Alias:
		return b.doAlias(sw)
	case types.Struct:
		return b.doStruct(sw)
	default:
		return b.unknown(sw)
	}
}

// GormFromGoNamer finds the gorm name of a type (and its package, and
// the package path) from its Go name.
type GormFromGoNamer interface {
	GoNameToGormName(name types.Name) types.Name
}

type GormLocator interface {
	GormTypeFor(t *types.Type) (*types.Type, error)
	GoTypeForName(name types.Name) *types.Type
	CastTypeName(name types.Name) string
}

type gormLocator struct {
	namer    GormFromGoNamer
	tracker  namer.ImportTracker
	universe types.Universe

	localGoPackage string
}

// CastTypeName returns the cast type name of a Go type
// TODO: delegate to a new localgo namer?
func (p gormLocator) CastTypeName(name types.Name) string {
	if name.Package == p.localGoPackage {
		return name.Name
	}
	return name.String()
}

func (p gormLocator) GoTypeForName(name types.Name) *types.Type {
	if len(name.Package) == 0 {
		name.Package = p.localGoPackage
	}
	return p.universe.Type(name)
}

// GormTypeFor locates a Gorm type for the provided Go type (if possible).
func (p gormLocator) GormTypeFor(t *types.Type) (*types.Type, error) {
	switch {
	// we've already converted the type, or it's a map
	case t.Kind == types.Gorm || t.Kind == types.Map:
		p.tracker.AddType(t)
		return t, nil
	}
	// it's a fundamental type
	if t, ok := isFundamentalGormType(t); ok {
		p.tracker.AddType(t)
		return t, nil
	}
	// it's a message
	if t.Kind == types.Struct || isOptionalAlias(t) {
		t := &types.Type{
			Name: p.namer.GoNameToGormName(t.Name),
			Kind: types.Gorm,

			CommentLines: t.CommentLines,
		}
		p.tracker.AddType(t)
		return t, nil
	}
	return nil, errUnrecognizedType
}

type bodyGen struct {
	locator        GormLocator
	localPackage   types.Name
	omitGogo       bool
	omitFieldTypes map[types.Name]struct{}

	t *types.Type
}

func (b bodyGen) unknown(sw *generator.SnippetWriter) error {
	return fmt.Errorf("not sure how to generate: %#v", b.t)
}

func (b bodyGen) doAlias(sw *generator.SnippetWriter) error {
	if !isOptionalAlias(b.t) {
		return nil
	}

	var kind string
	switch b.t.Underlying.Kind {
	case types.Map:
		kind = "map"
	default:
		kind = "slice"
	}
	optional := &types.Type{
		Name: b.t.Name,
		Kind: types.Struct,

		CommentLines:              b.t.CommentLines,
		SecondClosestCommentLines: b.t.SecondClosestCommentLines,
		Members: []types.Member{
			{
				Name:         "Items",
				CommentLines: []string{fmt.Sprintf("items, if empty, will result in an empty %s\n", kind)},
				Type:         b.t.Underlying,
			},
		},
	}
	nested := b
	nested.t = optional
	return nested.doStruct(sw)
}

func (b bodyGen) doStruct(sw *generator.SnippetWriter) error {
	if len(b.t.Name.Name) == 0 {
		return nil
	}
	if namer.IsPrivateGoName(b.t.Name.Name) {
		return nil
	}

	var alias *types.Type
	var fields []gormField
	if alias == nil {
		alias = b.t
	}

	// If we don't explicitly embed anything, generate fields by traversing fields.
	if fields == nil {
		memberFields, err := membersToFields(b.locator, alias, b.localPackage, b.omitFieldTypes)
		if err != nil {
			return fmt.Errorf("type %v cannot be converted to gorm: %v", b.t, err)
		}
		fields = memberFields
	}

	out := sw.Out()
	sw.Dof(`type XX_$.Name.Name$ struct {`, b.t)

	// primary key
	var pkField *gormField
	for i, field := range fields {
		if !extractFieldBoolTagOrDie(tagEnable, field.CommentLines) {
			continue
		}
		if field.primaryKey {
			pkField = &fields[i]
		}
		genComment(out, field.CommentLines, "  ")
		fmt.Fprintf(out, "  ")
		sw.Do("$.Name$ $.Type|local$ "+fmt.Sprintf("`%s`", field.toTagString()), field)
		fmt.Fprintf(out, "\n")
		if i != len(fields)-1 {
			fmt.Fprintf(out, "\n")
		}
	}
	sw.Doln("}")

	if pkField == nil {
		return fmt.Errorf("type %v missing field for primaryKey", b.t)
	}

	// generate storage struct
	sw.Dof(`type $.Name.Name$Storage struct {`, b.t)
	sw.Doln("*gorm.DB")
	sw.Doln("}")
	sw.Doln("")

	// generate New function for storage
	sw.Dof(`func New$.Name.Name$Storage(db *gorm.DB) *$.Name.Name$Storage {`, b.t)
	sw.Dof(`return &$.Name.Name$Storage{DB: db}`, b.t)
	sw.Doln("}")
	sw.Doln("")

	/*
		func (m TestStorage) PrimaryKey() (string, interface{}, bool) {
			return "uid", m.Uid, m.Uid == ""
		}
	*/
	// generate primary key
	sw.Dof(`func (s *$.Name.Name$Storage) PrimaryKey(m $.Name.Name$) (string, any, bool) {`, b.t)
	switch pkField.Type.Name.Name {
	case "string":
		sw.Dof(`return "$.GormName$", m.$.Name$, m.$.Name$ == ""`, pkField)
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		sw.Dof(`return "$.GormName$", m.$.Name$, m.$.Name$ == 0`, pkField)
	default:
		return fmt.Errorf("type %v invalid type for primaryKey", b.t)
	}
	sw.Doln("}")
	sw.Doln("")

	// generate CURD codes
	/*
		func (s *ProductStorage) extractClauses(ctx context.Context, m *Product) []clause.Expression {

			exprs := make([]clause.Expression, 0)
			if m.ID != 0 {
				exprs = append(exprs, Cond().Op(ParseOp(m.ID)).Build("id", m.ID))
			}

			//if m.Labels != nil {
			//	for k, v := range m.Labels {
			//		exprs = append(exprs, datatypes.JSONQuery("labels").Build())
			//	}
			//}

			return exprs
		}
	*/
	sw.Dof(`func (s *$.Name.Name$Storage) extractClauses(ctx context.Context, m *$.Name.Name$) []clause.Expression {`, b.t)
	sw.Doln("")
	sw.Doln(`exprs := make([]clause.Expression, 0)`)
	for _, field := range fields {
		if field.Type.Key != nil && field.Type.Elem != nil {
			sw.Dof(`if m.$.Name$ != nil {`, field)
			sw.Dof(`for k, v := range m.$.Name$ {`, field)
			sw.Dof(`exprs = append(exprs, datatypes.JSONQuery("$.GormName$").Equals(v, k))`, field)
			sw.Doln("}")
			sw.Doln("}")
		} else if field.Type.Elem != nil {
			sw.Dof(`if m.$.Name$ != nil {`, field)
			sw.Dof(`for _, item := range m.$.Name$ {`, field)
			if field.Type.Elem.Kind == "Builtin" {
				sw.Dof(`exprs = append(exprs, datatypes.JSONQuery("$.GormName$").HasKey(item))`, field)
			} else {
				sw.Dof(`for k, v := range dao.FieldPatch(item) {`, field)
				sw.Dof(`exprs = append(exprs, datatypes.JSONQuery("$.GormName$").Equals(v, strings.Split(k, ".")...))`, field)
				sw.Doln("}")
			}
			sw.Doln("}")
			sw.Doln("}")
		} else {
			switch field.Type.Name.Name {
			case "int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float", "double":
				if field.Type.Underlying != nil {
					sw.Dof(`if m.$.Name$ != nil {`, field)
					sw.Dof(`exprs = append(exprs, dao.Cond().Build("$.GormName$", *m.$.Name$))`, field)
					sw.Doln("}")
				} else {
					sw.Dof(`if m.$.Name$ != 0 {`, field)
					sw.Dof(`exprs = append(exprs, dao.Cond().Build("$.GormName$", m.$.Name$))`, field)
					sw.Doln("}")
				}
			case "string":
				if field.Type.Underlying != nil {
					sw.Dof(`if m.$.Name$ != nil {`, field)
					sw.Dof(`exprs = append(exprs, dao.Cond().Op(dao.ParseOp(*m.$.Name$)).Build("$.GormName$", *m.$.Name$))`, field)
					sw.Doln("}")
				} else {
					sw.Dof(`if m.$.Name$ != "" {`, field)
					sw.Dof(`exprs = append(exprs, dao.Cond().Op(dao.ParseOp(m.$.Name$)).Build("$.GormName$", m.$.Name$))`, field)
					sw.Doln("}")
				}
			case "bool":
				if field.Type.Underlying != nil {
					sw.Dof(`if m.$.Name$ != nil {`, field)
					sw.Dof(`exprs = append(exprs, dao.Cond().Build("$.GormName$", *m.$.Name$))`, field)
					sw.Doln("}")
				}
			}
		}
	}
	sw.Doln("")
	sw.Doln("return exprs")
	sw.Doln("}")
	sw.Doln("")

	return nil
}

type gormField struct {
	LocalPackage types.Name

	Name          string
	GormName      string
	Type          *types.Type
	Serializer    string
	Size          int
	primaryKey    bool
	Unique        bool
	Default       string
	Precision     bool
	Scale         bool
	Nullable      bool
	AutoIncrement bool
	Index         bool
	UniqueIndex   bool
	Extras        map[string]string

	CommentLines []string
}

func (f gormField) toTagString() string {
	items := []string{}
	items = append(items, "column:"+f.GormName)
	if f.Serializer != "" {
		items = append(items, "serializer:"+f.Serializer)
	}
	if f.primaryKey {
		items = append(items, tagPrimaryKey)
	}

	return "gorm:" + fmt.Sprintf(`"%s"`, strings.Join(items, ";"))
}

var (
	errUnrecognizedType = fmt.Errorf("did not recognize the provided type")
)

func isFundamentalGormType(t *types.Type) (*types.Type, bool) {
	// TODO: when we enable proto3, also include other fundamental types in the google.protobuf package
	// switch {
	// case t.Kind == types.Struct && t.Name == types.Name{Package: "time", Name: "Time"}:
	// 	return &types.Type{
	// 		Kind: types.Gorm,
	// 		Name: types.Name{Path: "google/protobuf/timestamp.proto", Package: "google.protobuf", Name: "Timestamp"},
	// 	}, true
	// }
	switch t.Kind {
	case types.Slice:
		if t.Elem.Name.Name == "byte" && len(t.Elem.Name.Package) == 0 {
			return &types.Type{Name: types.Name{Name: "bytes"}, Kind: types.Gorm}, true
		}
	case types.Builtin:
		switch t.Name.Name {
		case "string", "uint32", "int32", "uint64", "int64", "bool":
			return &types.Type{Name: types.Name{Name: t.Name.Name}, Kind: types.Gorm}, true
		case "int":
			return &types.Type{Name: types.Name{Name: "int64"}, Kind: types.Gorm}, true
		case "uint":
			return &types.Type{Name: types.Name{Name: "uint64"}, Kind: types.Gorm}, true
		case "float64", "float":
			return &types.Type{Name: types.Name{Name: "double"}, Kind: types.Gorm}, true
		case "float32":
			return &types.Type{Name: types.Name{Name: "float"}, Kind: types.Gorm}, true
		case "uintptr":
			return &types.Type{Name: types.Name{Name: "uint64"}, Kind: types.Gorm}, true
		}
		// TODO: complex?
	}
	return t, false
}

func memberTypeToGormField(locator GormLocator, field *gormField, t *types.Type, m *types.Member) error {
	var err error
	switch t.Kind {
	case types.Gorm:
		field.Type, err = locator.GormTypeFor(t)
	case types.Builtin:
		field.Type, err = locator.GormTypeFor(t)
	case types.Map:
		valueField := &gormField{}
		if err := memberTypeToGormField(locator, valueField, t.Elem, m); err != nil {
			return err
		}
		keyField := &gormField{}
		if err := memberTypeToGormField(locator, keyField, t.Key, m); err != nil {
			return err
		}
		// All other protobuf types have kind types.Gorm, so setting types.Map
		// here would be very misleading.
		field.Type = &types.Type{
			Kind: types.Gorm,
			Key:  keyField.Type,
			Elem: valueField.Type,
		}
		field.Serializer = "json"
	case types.Pointer:
		if err := memberTypeToGormField(locator, field, t.Elem, m); err != nil {
			return err
		}
		field.Type.Underlying = t.Elem
		field.Nullable = true
	case types.Alias:
		if isOptionalAlias(t) {
			field.Type, err = locator.GormTypeFor(t)
			field.Nullable = true
			field.Serializer = "json"
		} else {
			if err := memberTypeToGormField(locator, field, t.Underlying, m); err != nil {
				log.Warnf("failed to alias: %s %s: err %v", t.Name, t.Underlying.Name, err)
				return err
			}
		}
	case types.Slice:
		if t.Elem.Name.Name == "byte" && len(t.Elem.Name.Package) == 0 {
			field.Type = &types.Type{Name: types.Name{Name: "bytes"}, Kind: types.Gorm}
			return nil
		}
		if err := memberTypeToGormField(locator, field, t.Elem, m); err != nil {
			return err
		}
		field.Type.Elem = t.Elem
		field.Serializer = "json"
	case types.Struct:
		if len(t.Name.Name) == 0 {
			return errUnrecognizedType
		}
		field.Type, err = locator.GormTypeFor(t)
		field.Nullable = false
		field.Serializer = "json"
	default:
		return errUnrecognizedType
	}
	field.Name = m.Name
	return err
}

// gormTagToField extracts information from an existing gorm tag
func gormTagToField(tag string, field *gormField, m types.Member, t *types.Type, localPackage types.Name) error {
	if len(tag) == 0 || tag == "-" {
		return nil
	}

	// https://gorm.io/docs/models.html#Fields-Tags
	// gorm:"column:name;index;serializer:json"
	tMap := make(map[string]string)
	parts := strings.Split(tag, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) == 0 {
			continue
		}
		var k, v string
		if !strings.Contains(part, ":") {
			k = part
		} else {
			itemParts := strings.SplitN(part, ":", 2)
			k = itemParts[0]
			if len(itemParts) > 1 {
				v = itemParts[1]
			}
		}
		tMap[k] = v
	}

	field.Name = m.Name
	field.GormName = tMap["column"]
	field.Serializer = tMap["serializer"]
	if _, ok := tMap[tagPrimaryKey]; ok {
		field.primaryKey = true
	}

	return nil
}

func membersToFields(locator GormLocator, t *types.Type, localPackage types.Name, omitFieldTypes map[types.Name]struct{}) ([]gormField, error) {
	fields := []gormField{}

	for _, m := range t.Members {
		if namer.IsPrivateGoName(m.Name) {
			// skip private fields
			continue
		}
		if !extractFieldBoolTagOrDie(tagEnable, m.CommentLines) {
			continue
		}
		if _, ok := omitFieldTypes[types.Name{Name: m.Type.Name.Name, Package: m.Type.Name.Package}]; ok {
			continue
		}
		tags := reflect.StructTag(m.Tags)
		field := gormField{
			LocalPackage: localPackage,
			//Tag:          -1,
			Extras: make(map[string]string),
		}

		markers := types.ExtractCommentTags("+", m.CommentLines)
		if v := markers[tagPrimaryKey]; v != nil {
			field.primaryKey = true
		}

		gormTag := tags.Get("gorm")
		if gormTag == "-" {
			continue
		}

		if err := gormTagToField(gormTag, &field, m, t, localPackage); err != nil {
			return nil, err
		}

		// extract information from JSON field tag
		if tag := tags.Get("json"); len(tag) > 0 {
			parts := strings.Split(tag, ",")
			if len(field.Name) == 0 && len(parts[0]) != 0 {
				field.Name = parts[0]
			}
			if field.Name == "-" {
				continue
			}
			i := 0
			length := len(field.Name)
			buf := &bytes.Buffer{}
			for i < length {
				c := field.Name[i]
				if c == '.' || c == '-' {
					c = '_'
				}
				i += 1
				buf.WriteByte(c)
			}
			field.GormName = buf.String()
		}

		if field.Type == nil {
			if err := memberTypeToGormField(locator, &field, m.Type, &m); err != nil {
				return nil, fmt.Errorf("unable to embed type %q as field %q in %q: %v", m.Type, field.Name, t.Name, err)
			}
		}
		if len(field.Name) == 0 {
			field.Name = m.Name
		}

		field.CommentLines = m.CommentLines
		fields = append(fields, field)
	}

	return fields, nil
}

func genComment(out io.Writer, lines []string, indent string) {
	for {
		l := len(lines)
		if l == 0 || len(lines[l-1]) != 0 {
			break
		}
		lines = lines[:l-1]
	}
	for _, c := range lines {
		if len(c) == 0 {
			fmt.Fprintf(out, "%s//\n", indent) // avoid trailing whitespace
			continue
		}
		fmt.Fprintf(out, "%s// %s\n", indent, c)
	}
}

func formatGormFile(source []byte) ([]byte, error) {
	// TODO; Is there any protobuf formatter?
	return source, nil
}

func assembleGormFile(w io.Writer, f *generator.File) {
	w.Write(f.Header)

	if f.Vars.Len() > 0 {
		fmt.Fprintf(w, "%s\n", f.Vars.String())
	}

	//f.Imports["gorm.io/gorm"] = struct{}{}

	if len(f.Imports) > 0 {
		imports := []string{}
		for i := range f.Imports {
			imports = append(imports, i)
		}
		sort.Strings(imports)
		fmt.Fprintf(w, "import (\n")
		for _, s := range imports {
			fmt.Fprintf(w, `"%s"`, s)
			fmt.Fprint(w, "\n")
		}
		fmt.Fprintf(w, ")\n")
	}

	w.Write(f.Body.Bytes())
}

func NewGormFile() *generator.DefaultFileType {
	return &generator.DefaultFileType{
		Format:   formatGormFile,
		Assemble: assembleGormFile,
	}
}
