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
	"github.com/vine-io/gogogen/gogenerator/namer"
	"github.com/vine-io/gogogen/gogenerator/types"
)

type ImportTracker struct {
	namer.DefaultImportTracker
}

func NewImportTracker(local types.Name, typesToAdd ...*types.Type) *ImportTracker {
	tracker := namer.NewDefaultImportTracker(local)
	tracker.IsInvalidType = func(t *types.Type) bool { return t.Kind != types.Gorm }
	tracker.LocalName = func(name types.Name) string { return name.Package }
	tracker.PrintImport = func(path, name string) string { return path }

	tracker.AddTypes(typesToAdd...)
	return &ImportTracker{
		DefaultImportTracker: tracker,
	}
}
