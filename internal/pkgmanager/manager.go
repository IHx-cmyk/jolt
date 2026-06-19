package pkgmanager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

const installDir = ".jolt/packages"

type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	Main            string            `json:"main"`
	Author          string            `json:"author"`
	License         string            `json:"license"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func Add(pkgName string) error {
	pkgFile := "package.json"
	var pkg PackageJSON

	if data, err := ioutil.ReadFile(pkgFile); err == nil {
		if err := json.Unmarshal(data, &pkg); err != nil {
			return fmt.Errorf("gagal parse package.json: %w", err)
		}
	} else {
		wd, _ := os.Getwd()
		pkg = PackageJSON{
			Name:            filepath.Base(wd),
			Version:         "1.0.0",
			Description:     "",
			Main:            "index.js",
			Author:          "",
			License:         "MIT",
			Dependencies:    make(map[string]string),
			DevDependencies: make(map[string]string),
		}
	}
	if pkg.Dependencies == nil {
		pkg.Dependencies = make(map[string]string)
	}

	startTime := time.Now()

	dotChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	s := spinner.New(dotChars, 80*time.Millisecond)
	s.Suffix = "  Mengambil info package..."
	s.Color("cyan")
	s.Start()

	info, err := fetchPackageInfo(pkgName)
	if err != nil {
		s.Stop()
		return err
	}
	latestVersion := info["dist-tags"].(map[string]interface{})["latest"].(string)
	versionInfo := info["versions"].(map[string]interface{})[latestVersion].(map[string]interface{})
	tarballURL := versionInfo["dist"].(map[string]interface{})["tarball"].(string)

	s.Suffix = fmt.Sprintf("  Download %s@%s... (%s)", pkgName, latestVersion, time.Since(startTime).Round(time.Second))
	s.Restart()

	targetDir := filepath.Join(installDir, pkgName)
	os.MkdirAll(targetDir, 0755)

	if err := downloadAndExtractLight(tarballURL, targetDir); err != nil {
		s.Stop()
		return err
	}

	s.Suffix = fmt.Sprintf("  Menyimpan konfigurasi... (%s)", time.Since(startTime).Round(time.Second))
	s.Restart()
	time.Sleep(500 * time.Millisecond)

	pkg.Dependencies[pkgName] = latestVersion
	data, _ := json.MarshalIndent(pkg, "", "  ")
	ioutil.WriteFile(pkgFile, data, 0644)

	s.Stop()
	elapsed := time.Since(startTime).Round(time.Second)
	color.Green("✅ %s@%s berhasil diinstall (selesai dalam %s)", pkgName, latestVersion, elapsed)
	color.Cyan("   📦 Lokasi: .jolt/packages/%s", pkgName)
	return nil
}