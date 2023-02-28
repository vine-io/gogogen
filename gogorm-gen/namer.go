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
	"reflect"
	"strings"

	"github.com/vine-io/gogogen/gogenerator/generator"
	"github.com/vine-io/gogogen/gogenerator/namer"
	"github.com/vine-io/gogogen/gogenerator/types"
)

type localNamer struct {
	localPackage types.Name
}

func (n localNamer) Name(t *types.Type) string {
	if t.Key != nil && t.Elem != nil {
		return fmt.Sprintf("map[%s]%s", n.Name(t.Key), n.Name(t.Elem))
	}
	if len(n.localPackage.Package) != 0 && n.localPackage.Package == t.Name.Package {
		return t.Name.Name
	}
	return t.Name.String()
}

type gormNamer struct {
	packages       []*gormPackage
	packagesByPath map[string]*gormPackage
}

func NewGormNamer() *gormNamer {
	return &gormNamer{
		packagesByPath: make(map[string]*gormPackage),
	}
}

func (n *gormNamer) Name(t *types.Type) string {
	if t.Kind == types.Map {
		return fmt.Sprintf("map[%s]%s", n.Name(t.Key), n.Name(t.Elem))
	}
	return t.Name.String()
}

func (n *gormNamer) List() []generator.Package {
	packages := make([]generator.Package, 0, len(n.packages))
	for i := range n.packages {
		packages = append(packages, n.packages[i])
	}
	return packages
}

func (n *gormNamer) Add(p *gormPackage) {
	if _, ok := n.packagesByPath[p.PackagePath]; !ok {
		n.packagesByPath[p.PackagePath] = p
		n.packages = append(n.packages, p)
	}
}

func (n *gormNamer) GoNameToGormName(name types.Name) types.Name {
	if p, ok := n.packagesByPath[name.Package]; ok {
		return types.Name{
			Name:    name.Name,
			Package: p.PackageName,
			//Path:    p.ImportPath(),
		}
	}
	for _, p := range n.packages {
		if _, ok := p.FilterTypes[name]; ok {
			return types.Name{
				Name:    name.Name,
				Package: p.PackageName,
				//Path:    p.ImportPath(),
			}
		}
	}
	return types.Name{Name: name.Name}
}

func gormSafePackage(name string) string {
	pkg := strings.Replace(name, "/", ".", -1)
	return strings.Replace(pkg, "-", "_", -1)
}

type typeNameSet map[types.Name]*gormPackage

// assignGoTypeToGormPackage looks for Go types that are referenced by a type in a package.
func assignGoTypeToGormPackage(p *gormPackage, t *types.Type, local, global typeNameSet, optional map[types.Name]struct{}) {
	newT, isProto := isFundamentalGormType(t)
	if isProto {
		t = newT
	}
	if otherP, ok := global[t.Name]; ok {
		if _, ok := local[t.Name]; !ok {
			p.Imports.AddType(&types.Type{
				Kind: types.Gorm,
				Name: otherP.GormTypeName(),
			})
		}
		return
	}
	if t.Name.Package == p.PackagePath {
		// Associate types only to their own package
		global[t.Name] = p
	}
	if _, ok := local[t.Name]; ok {
		return
	}
	// don't recurse into existing gorm types
	if isProto {
		p.Imports.AddType(t)
		return
	}

	local[t.Name] = p
	for _, m := range t.Members {
		if namer.IsPrivateGoName(m.Name) {
			continue
		}
		field := &gormField{}
		tag := reflect.StructTag(m.Tags).Get("gorm")
		if tag == "-" {
			continue
		}
		if err := gormTagToField(tag, field, m, t, p.GormTypeName()); err == nil && field.Type != nil {
			assignGoTypeToGormPackage(p, field.Type, local, global, optional)
			continue
		}
		assignGoTypeToGormPackage(p, m.Type, local, global, optional)
	}
	// TODO: should methods be walked?
	if t.Elem != nil {
		assignGoTypeToGormPackage(p, t.Elem, local, global, optional)
	}
	if t.Key != nil {
		assignGoTypeToGormPackage(p, t.Key, local, global, optional)
	}
	if t.Underlying != nil {
		if t.Kind == types.Alias && isOptionalAlias(t) {
			optional[t.Name] = struct{}{}
		}
		assignGoTypeToGormPackage(p, t.Underlying, local, global, optional)
	}
}

// isTypeApplicableToGorm checks to see if a type is relevant for protobuf processing.
// Currently, it filters out functions and private types.
func isTypeApplicableToGorm(t *types.Type) bool {
	// skip functions -- we don't care about them for protobuf
	if t.Kind == types.Func || (t.Kind == types.DeclarationOf && t.Underlying.Kind == types.Func) {
		return false
	}
	// skip private types
	if namer.IsPrivateGoName(t.Name.Name) {
		return false
	}

	return true
}

func (n *gormNamer) AssignTypesToPackages(c *generator.Context) error {
	global := make(typeNameSet)
	for _, p := range n.packages {
		local := make(typeNameSet)
		optional := make(map[types.Name]struct{})
		p.Imports = NewImportTracker(p.GormTypeName())
		for _, t := range c.Order {
			if t.Name.Package != p.PackagePath {
				continue
			}
			if !isTypeApplicableToGorm(t) {
				// skip types that we don't care about, like functions
				continue
			}
			assignGoTypeToGormPackage(p, t, local, global, optional)
		}
		p.FilterTypes = make(map[types.Name]struct{})
		p.LocalNames = make(map[string]struct{})
		p.OptionalTypeNames = make(map[string]struct{})
		for k, v := range local {
			if v == p {
				p.FilterTypes[k] = struct{}{}
				p.LocalNames[k.Name] = struct{}{}
				if _, ok := optional[k]; ok {
					p.OptionalTypeNames[k.Name] = struct{}{}
				}
			}
		}
	}
	return nil
}
