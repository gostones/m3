package internal

import (
	"os"
	"path/filepath"
	"strings"
)

// PathJoinList is the inverse operation of filepath.SplitList
func PathJoinList(p []string) string {
	return strings.Join(p, string(filepath.ListSeparator))
}

// AddPath adds list of pathes the PATH env per OS convention
func AddPath(p []string) {
	env := os.Getenv("PATH")
	pl := filepath.SplitList(env)
	for _, i := range pl {
		p = append(p, i)
	}
	os.Setenv("PATH", strings.Join(p, string(filepath.ListSeparator)))
}
