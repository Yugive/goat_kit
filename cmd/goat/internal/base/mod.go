package base

import (
	"bufio"
	"bytes"
	"golang.org/x/mod/modfile"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//Get go module path
func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

//Get go module version
func ModuleVersion(path string) (string, error) {
	stdout := &bytes.Buffer{}
	fd := exec.Command("go", "mod", "graph")
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	if err := fd.Run(); err != nil {
		return "", err
	}
	rd := bufio.NewReader(stdout)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			return "", err
		}
		str := string(line)
		i := strings.Index(str, "@")
		if strings.Contains(str, path+"@") && i != -1 {
			return path + str[i:], nil
		}
	}
}

//Get goat mod
func GoatMod() string {
	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
	cachePath := strings.Trim(string(cacheOut), "\n")
	pathOut, _ := exec.Command("go", "env", "GOPATH").Output()
	goPath := strings.Trim(string(pathOut), "\n")
	if cachePath == "" {
		cachePath = filepath.Join(goPath, "mod", "pkg")
	}
	if path, err := ModuleVersion("github.com/ngyugive/go-web-framework"); err == nil {
		return filepath.Join(cachePath, path)
	}
	return filepath.Join(goPath, "src", "github.com", "ngyugive", "go-web-framework")
}
