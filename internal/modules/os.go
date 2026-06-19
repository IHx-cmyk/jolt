package modules

import (
	"os"
	"github.com/dop251/goja"
)

func RegisterOS(vm *goja.Runtime) {
	osModule := vm.NewObject()
	osModule.Set("getenv", func(key string) string {
		return os.Getenv(key)
	})
	osModule.Set("args", func() []string {
		return os.Args[1:]
	})
	vm.Set("os", osModule)
}