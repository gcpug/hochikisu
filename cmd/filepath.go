package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsFullPath(exe string) bool {
	if len(exe) < 1 {
		return false
	}

	if exe[0:1] == "/" {
		return true
	}

	return false
}

func ReadFile(path string) ([]byte, error) {
	var fp string
	if IsFullPath(path) {
		fp = path
	} else {
		exe, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed os.Executable(). err=%+v", err)
		}
		pwd := filepath.Dir(exe)
		fp = fmt.Sprintf("%s/%s", pwd, path)
	}

	body, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, fmt.Errorf("failed read file. filepath=%s, err=%+v", fp, err)
	}
	return body, nil
}
