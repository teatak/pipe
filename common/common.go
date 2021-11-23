package common

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var BinDir string

func init() {
	rand.Seed(time.Now().UnixNano())
	bin, _ := os.Executable()
	realPath, err := os.Readlink(bin)
	if err == nil {
		bin = realPath
	}
	if filepath.Base(bin) == Name {
		BinDir = filepath.Dir(bin)
	} else {
		BinDir, _ = os.Getwd()
	}
}

func MakeFile(path string) *os.File {
	dir := filepath.Dir(path)
	if !IsExist(dir) {
		os.MkdirAll(dir, 0755)
	}
	file, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return file
}

func IsExist(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}
