package aspera

import (
	"fmt"
	"runtime"
)

const (
	version = "1.1.1"
	commit  = "52a85ef"
	prefix  = "https://download.asperasoft.com/download/sw/sdk/transfer"
)

func GetSDKDownloadURL() (url string, platform string, err error) {
	platform, ext, err := GetSDKAttributes(runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s/%s/%s-%s.%s", prefix, version, platform, commit, ext)
	return
}

func GetSDKAttributes(os, arch string) (platform string, ext string, err error) {
	platforms := map[string][]string{
		"darwin":  {"amd64"},
		"linux":   {"amd64", "ppc64le", "s390x"},
		"windows": {"amd64"},
		"aix":     {"ppc64"},
	}

	ext = "tar.gz"

	if archs, ok := platforms[os]; ok {
		for _, a := range archs {
			if a == arch {
				if os == "darwin" {
					os = "osx"
				}
				if os == "windows" {
					ext = "zip"
				}
				return fmt.Sprintf("%s-%s", os, arch), ext, nil
			}
		}
	}
	return "", "", fmt.Errorf("unsupported platform: %s-%s", os, arch)

}
