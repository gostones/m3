package internal

import (
	"log"
	"os"
)

var Stdlog, Errlog *log.Logger

func init() {
	Stdlog = log.New(os.Stdout, "", 0)
	Errlog = log.New(os.Stderr, "", 0)
}

func DumpEnv() {
	Stdlog.Println("dump env ...")

	for _, nv := range os.Environ() {
		Stdlog.Println(nv)
	}
}
