package engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/dop251/goja"
	"jolt/internal/modules"
)

type Runtime struct {
	vm *goja.Runtime
}

func New(opts Options) *Runtime {
	vm := goja.New()
	modules.RegisterConsole(vm)
	modules.RegisterFS(vm)
	modules.RegisterNet(vm)
	modules.RegisterOS(vm)

	cache := make(map[string]goja.Value)

	vm.Set("require", func(moduleName string) goja.Value {
		if val, ok := cache[moduleName]; ok {
			return val
		}

		baseDir := ".jolt/packages"
		possiblePaths := []string{
			filepath.Join(baseDir, moduleName, "index.js"),
			filepath.Join(baseDir, moduleName, moduleName+".js"),
		}

		var mainFile string
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				mainFile = p
				break
			}
		}
		if mainFile == "" {
			panic(vm.NewGoError(fmt.Errorf("module %s tidak ditemukan", moduleName)))
		}

		code, err := ioutil.ReadFile(mainFile)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		val, err := vm.RunString(string(code))
		if err != nil {
			panic(vm.NewGoError(err))
		}
		cache[moduleName] = val
		return val
	})

	return &Runtime{vm: vm}
}

func (r *Runtime) RunFile(path string) error {
	code, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = r.vm.RunString(string(code))
	return err
}