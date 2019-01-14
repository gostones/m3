package internal

import (
	"log"
	"os"
)

var Stdlog, Errlog *log.Logger

var logger *log.Logger

func init() {
	Stdlog = log.New(os.Stdout, "", 0)
	Errlog = log.New(os.Stderr, "", 0)
	logger = Stdlog
}

func DumpEnv() {
	Stdlog.Println("dump env ...")

	for _, nv := range os.Environ() {
		Stdlog.Println(nv)
	}
}
