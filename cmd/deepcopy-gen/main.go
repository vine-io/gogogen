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

// deepcopy-gen is a tool for auto-generating DeepCopy functions.
//
// Given a list of input directories, it will generate functions that
// efficiently perform a full deep-copy each type. For any type that
// offers a `.DeepCopy()` method, it will simple call that. Otherwise it will
// use standard value assignment whenever possible. If that is not possible it
// will try to call its own generated copy function for the type, if the type is
// within the allowed root packages. Failing that, it will fall back on
// `conversion.Cloner.DeepCopy(val)` to make to copy. The resulting file will
// be stored in the same directory as the processed source package.
//
// Generation is governed by comment tags in the source. Any package may
// request DeepCopy generation by including a comment in the file-comments of
// one file, of the form:
//
//	// +gogo:deepcopy-gen=package
//
// DeepCopy functions can be generated for individual types, rather than the
// entire package by specified a comment on the type define of the form:
//
//	// +gogo:deepcopy-gen=true
//
// When generating for a whole package, individual type may opt out of
// DeepCopy generation by specified a comment on the of form:
//
//	// +gogo:deepcopy-gen=false
//
// Note that registration is a whole-package option, and is not available for
// individual types.
package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	deepcopy_gen "github.com/vine-io/gogogen/deepcopy-gen"
	"github.com/vine-io/gogogen/gogenerator/args"
	"github.com/vine-io/gogogen/util/log"

	utilbuild "github.com/vine-io/gogogen/util/build"
)

func main() {
	genericArgs, customArgs := deepcopy_gen.NewDefaults()

	// Override defaults.
	// TODO: move this out of deepcopy-gen
	genericArgs.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), utilbuild.BoilerplatePath())

	fs := pflag.NewFlagSet("deepcopy", pflag.ExitOnError)
	genericArgs.AddFlags(fs)
	customArgs.AddFlags(fs)
	if err := fs.Parse(os.Args); err != nil {
		log.Fatalf("Error: %v", err)
	}

	if err := deepcopy_gen.Validate(genericArgs); err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Run it.
	if err := genericArgs.Execute(
		deepcopy_gen.NameSystems(),
		deepcopy_gen.DefaultNameSystem(),
		deepcopy_gen.Package,
	); err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Infof("Completed successfully.")
}
