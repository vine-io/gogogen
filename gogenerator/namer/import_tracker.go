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

package namer

import (
	"github.com/vine-io/gogogen/gogenerator/types"
)

// ImportTracker may be passed to a namer.RawNamer, to tracker the imports needed
// for the type it names.
//
// TODO: pay attention to the package name (instead of renaming every package).
type DefaultImportTracker struct {
	pathToName map[string]string
	// forbidden names are in here. (e.g. "go" is a directory in which
	// there is code, but "go" is not a legal name for a package, so we put
	// it here to prevent us from naming any package "go")
	nameToPath map[string]string
	local      types.Name

	// Returns true if a given types an invalid type and should be ignored.
	IsInvalidType func(*types.Type) bool
	// Returns the final local name for the given name
	LocalName func(types.Name) string
	// Returns the "import" line for a given (path, name).
	PrintImport func(string, string) string
}

func NewDefaultImportTracker(local types.Name) DefaultImportTracker {
	return DefaultImportTracker{
		pathToName: map[string]string{},
		nameToPath: map[string]string{},
		local:      local,
	}
}

func (tracker *DefaultImportTracker) AddTypes(types ...*types.Type) {
	for _, t := range types {
		tracker.AddType(t)
	}
}

func (tracker *DefaultImportTracker) AddType(t *types.Type) {
	if tracker.local.Package == t.Name.Package {
		return
	}

	if tracker.IsInvalidType(t) {
		if t.Kind == types.Builtin {
			return
		}
		if _, ok := tracker.nameToPath[t.Name.Package]; !ok {
			tracker.nameToPath[t.Name.Package] = ""
		}
		return
	}

	if len(t.Name.Package) == 0 {
		return
	}
	path := t.Name.Path
	if len(path) == 0 {
		path = t.Name.Package
	}
	if _, ok := tracker.pathToName[path]; ok {
		return
	}
	name := tracker.LocalName(t.Name)
	tracker.nameToPath[name] = path
	tracker.pathToName[path] = name
}

func (tracker *DefaultImportTracker) ImportLines() map[string]string {
	out := map[string]string{}
	for path, name := range tracker.pathToName {
		out[path] = name
	}
	return out
}

// LocalNameOf returns the name you would use to refer to the package at the
// specified path within the way body of a file.
func (tracker *DefaultImportTracker) LocalNameOf(path string) string {
	return tracker.pathToName[path]
}

// PathOf returns the path that a given localName is referring to within the
// body of a file.
func (tracker *DefaultImportTracker) PathOf(localPath string) (string, bool) {
	name, ok := tracker.nameToPath[localPath]
	return name, ok
}
