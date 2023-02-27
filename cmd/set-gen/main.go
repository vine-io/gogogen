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

// set-gen is an example usage of gengo
//
// Struct in the input directories with the below line in their comments will
// have sets generated for them.
// // +gogo:genset
//
// Any builtin type referenced anywhere in the input directories will have a
// set generated for it.
package main

import (
	"os"
	"path/filepath"

	"github.com/vine-io/gogogen/gogenerator/args"
	set_gen "github.com/vine-io/gogogen/set-gen"
	utilbuild "github.com/vine-io/gogogen/util/build"
	"github.com/vine-io/gogogen/util/log"
)

func main() {
	arguments := args.Default()

	// Override defaults.
	arguments.GoHeaderFilePath = filepath.Join(args.DefaultSourceTree(), utilbuild.BoilerplatePath())
	arguments.InputDirs = []string{"github.com/vine-io/gogogen/util/sets/types"}
	arguments.OutputPackagePath = "github.com/vine-io/gogogen/util/sets"

	if err := arguments.Execute(
		set_gen.NameSystems(),
		set_gen.DefaultNameSystem(),
		set_gen.Packages,
	); err != nil {
		log.Errorf("Error: %v", err)
		os.Exit(1)
	}
	log.Infof("Completed successfully.")
}
