package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PathJoinList is the inverse operation of filepath.SplitList
func PathJoinList(p []string) string {
	return strings.Join(p, string(filepath.ListSeparator))
}

// AddPath adds list of pathes the PATH env per OS convention
func AddPath(env string, p []string) string {
	pl := filepath.SplitList(env)
	for _, i := range pl {
		p = append(p, i)
	}
	return strings.Join(p, string(filepath.ListSeparator))
}

// DefaultEnviron returns required env slice for running core services
func DefaultEnviron(base string) []string {
	home := os.Getenv("HOME")
	if home == "" {
		home = filepath.Dir(base)
	}
	//TODO template?
	el := []string{
		fmt.Sprintf("HOME=%v", home),
		fmt.Sprintf("DHNT_BASE=%v", base),
		fmt.Sprintf("GOPATH=%v/go", base),
		fmt.Sprintf("IPFS_PATH=%v/home/ipfs", base),
		fmt.Sprintf("GOGS_WORK_DIR=%v/home/gogs", base),
		fmt.Sprintf("PATH=%v", AddPath(os.Getenv("PATH"), []string{
			fmt.Sprintf("%v/go/bin", base),
			fmt.Sprintf("%v/bin", base),
		})),
	}
	env := os.Environ()
	for _, e := range el {
		env = append(env, e)
	}
	return env
}

// SetDefaultEnviron sets required env for running core services
func SetDefaultEnviron(base string) {
	env := DefaultEnviron(base)
	for _, v := range env {
		nv := strings.SplitN(v, "=", 2)
		os.Setenv(nv[0], nv[1])
	}
}

// GetIntEnv returns int env or default
func GetIntEnv(env string, i int) int {
	if p := os.Getenv(env); p != "" {
		if i, err := strconv.Atoi(p); err == nil {
			return i
		}
	}
	return i
}
