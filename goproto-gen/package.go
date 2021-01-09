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
	"fmt"
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/lack-io/gogogen/gogenerator/generator"
	"github.com/lack-io/gogogen/gogenerator/types"
	"github.com/lack-io/gogogen/util/log"
)

func newProtobufPackage(packagePath, packageName string, generateAll bool, omitFieldTypes map[types.Name]struct{}) *protobufPackage {
	pkg := &protobufPackage{
		DefaultPackage: generator.DefaultPackage{
			// The protobuf package name (foo.bar.baz)
			PackageName: packageName,
			// A path segment relative to the GOPATH root (foo/bar/baz)
			PackagePath: packagePath,
			HeaderText: []byte(
				`
// This file was autogenerated by goproto-gen. Do not edit it manually!

`),
			PackageDocumentation: []byte(fmt.Sprintf(
				`// Package %s is an autogenerated protobuf IDL.
`, packageName)),
		},
		GenerateAll:    generateAll,
		OmitFieldTypes: omitFieldTypes,
	}
	pkg.FilterFunc = pkg.filterFunc
	pkg.GeneratorFunc = pkg.generatorFunc
	return pkg
}

// protobufPackage contains the protobuf implementation of Package.
type protobufPackage struct {
	generator.DefaultPackage

	// If true, this package has been vendored into our source tree and thus can
	// only be generated by changing the vendor tree.
	Vendored bool

	// If true, generate protobuf serializations for all public types.
	// If false, only generate protobuf serializations for structs that
	// request serialization.
	GenerateAll bool

	// A list of types to filter to; if not specified all types will be included.
	FilterTypes map[types.Name]struct{}

	// If true, omit any gogoprotobuf extensions not defined as types.
	OmitGogo bool

	// A list of field types that will be excluded from the output struct
	OmitFieldTypes map[types.Name]struct{}

	// A list of names that this package exports
	LocalNames map[string]struct{}

	// A list of type names in this package that will need marshaller rewriting
	// to remove synthetic protobuf fields.
	OptionalTypeNames map[string]struct{}

	// A list of struct tags to generate onto named struct fields
	StructTags map[string]map[string]string

	// An import tracker for this package
	Imports *ImportTracker
}

func (p *protobufPackage) Clean(outputBase string) error {
	for _, s := range []string{p.ImportPath(), p.OutputPath()} {
		if err := os.Remove(filepath.Join(outputBase, s)); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (p *protobufPackage) ProtoTypeName() types.Name {
	return types.Name{
		Name:    p.Path(),       // the go path "foo/bar/baz"
		Package: p.Name(),       // the protobuf package "foo.bar.baz"
		Path:    p.ImportPath(), // the path of the import to get the proto
	}
}

func (p *protobufPackage) filterFunc(c *generator.Context, t *types.Type) bool {
	switch t.Kind {
	case types.Func, types.Chan:
		return false
	case types.Struct:
		if t.Name.Name == "struct{}" {
			return false
		}
		// +gogo:genproto
		if !extractBoolTagOrDie(tagEnable, t.CommentLines) {
			return false
		}
	case types.Builtin:
		return false
	case types.Alias:
		if !isOptionalAlias(t) {
			return false
		}
	case types.Slice, types.Array, types.Map:
		return false
	case types.Pointer:
		return false
	}
	if _, ok := isFundamentalProtoType(t); ok {
		return false
	}
	_, ok := p.FilterTypes[t.Name]
	return ok
}

func (p *protobufPackage) HasGoType(name string) bool {
	_, ok := p.LocalNames[name]
	return ok
}

func (p *protobufPackage) OptionalTypeName(name string) bool {
	_, ok := p.OptionalTypeNames[name]
	return ok
}

func (p *protobufPackage) ExtractGeneratedType(t *ast.TypeSpec) bool {
	if !p.HasGoType(t.Name.Name) {
		return false
	}

	switch s := t.Type.(type) {
	case *ast.StructType:
		for i, f := range s.Fields.List {
			if len(f.Tag.Value) == 0 {
				continue
			}
			tag := strings.Trim(f.Tag.Value, "`")
			protobufTag := reflect.StructTag(tag).Get("protobuf")
			if len(protobufTag) == 0 {
				continue
			}
			if len(f.Names) > 1 {
				log.Warnf("WARNING: struct %s field %d %s: defined multiple names but single protobuf tag", t.Name.Name, i, f.Names[0].Name)
				// TODO hard error?
			}
			if p.StructTags == nil {
				p.StructTags = make(map[string]map[string]string)
			}
			m := p.StructTags[t.Name.Name]
			if m == nil {
				m = make(map[string]string)
				p.StructTags[t.Name.Name] = m
			}
			m[f.Names[0].Name] = tag
		}
	default:
		log.Warnf("WARNING: unexpected Go AST type definition: %#v", t)
	}

	return true
}

func (p *protobufPackage) generatorFunc(c *generator.Context) []generator.Generator {
	generators := []generator.Generator{}

	p.Imports.AddNullable()

	generators = append(generators, &genProtoIDL{
		DefaultGen: generator.DefaultGen{
			OptionalName: "generated",
		},
		localPackage:   types.Name{Package: p.PackageName, Path: p.PackagePath},
		localGoPackage: types.Name{Package: p.PackagePath, Name: p.GoPackageName()},
		imports:        p.Imports,
		generateAll:    p.GenerateAll,
		omitGogo:       p.OmitGogo,
		omitFieldTypes: p.OmitFieldTypes,
	})
	return generators
}

func (p *protobufPackage) GoPackageName() string {
	return filepath.Base(p.PackagePath)
}

func (p *protobufPackage) ImportPath() string {
	return filepath.Join(p.PackagePath, "generated.proto")
}

func (p *protobufPackage) OutputPath() string {
	return filepath.Join(p.PackagePath, "generated.pb.go")
}

var (
	_ = generator.Package(&protobufPackage{})
)
