package internal

import (
	"path/filepath"
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
