package pkgmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func fetchPackageInfo(pkg string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s", pkg)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("package tidak ditemukan (status %d)", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	return data, err
}