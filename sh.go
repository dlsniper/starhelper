// Copyright 2020 Florin Pățan
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Starhelper makes binding to user-defined Go functions easier by providing the convenience
// to register these functions and handle their execution.
package starhelper

import (
	"fmt"
	"runtime"

	"go.starlark.net/starlark"
)

// App holds our application execution context
type App struct {
	preDeclared starlark.StringDict
}

// Execute calls the script and runs our application
func (env *App) Execute(thread *starlark.Thread, path string) error {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		fmt.Printf("panic while running the starlark script: %v\n", err)
		for i := 3; ; i++ {
			pc, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			fname := "<unknown>"
			fn := runtime.FuncForPC(pc)
			if fn != nil {
				fname = fn.Name()
			}
			fmt.Printf("%s\n\tin %s:%d\n", fname, file, line)
		}
	}()

	// Discard the return value as we are not interested in it
	_, err := starlark.ExecFile(thread, path, nil, env.preDeclared)
	return err
}

// New generates a new execution context for a Starlark based application
func New(builtins starlark.StringDict) *App {
	return &App{preDeclared: builtins}
}
