package rp

import (
	"fmt"
	"github.com/dhnt/m3/internal/misc"
	"io/ioutil"
	"path/filepath"
)

var logger = misc.Stdlog

// checkOrCreateEskip returns config filename; creates one if it does not exist.
func checkOrCreateEskip(base, myid string) (string, error) {
	cf := filepath.Join(base, "etc/routes.eskip")
	if _, err := misc.CreateDir(filepath.Dir(cf)); err != nil {
		return "", err
	}
	logger.Println("Skipper routes file: ", cf)

	if misc.ExistRegularFile(cf) {
		return cf, nil
	}

	data := []byte(fmt.Sprintf(routesEskip, myid, myid, myid))
	if err := ioutil.WriteFile(cf, data, 0644); err != nil {
		return "", err
	}
	return cf, nil
}

// TODO Dataclients https://github.com/zalando/skipper/blob/master/docs/tutorials/development.md
var routesEskip = `
riot:
    Host("^riot.(home|%v)$")
    -> setRequestHeader("X-Skipper-Tag", "skipper")
    -> "http://localhost:8080";

matrix:
    Host("^matrix.(home|%v)(:8008)?$")
    -> setRequestHeader("X-Skipper-Tag", "skipper")
    -> "http://localhost:8008";

matrixFederation:
    Host("^matrix.(home|%v):8448$")
    -> setRequestHeader("X-Skipper-Tag", "skipper")
	-> "http://localhost:8448";

traefik: * -> setRequestHeader("X-Skipper-Tag", "skipper") -> "http://127.0.0.1:80"
`
