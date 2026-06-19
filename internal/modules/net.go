package modules

import (
	"io/ioutil"
	"net/http"
	"github.com/dop251/goja"
)

func RegisterNet(vm *goja.Runtime) {
	net := vm.NewObject()
	net.Set("fetch", func(url string) map[string]interface{} {
		resp, err := http.Get(url)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(vm.NewGoError(err))
		}
		return map[string]interface{}{
			"status": resp.StatusCode,
			"body":   string(body),
		}
	})
	vm.Set("net", net)
}