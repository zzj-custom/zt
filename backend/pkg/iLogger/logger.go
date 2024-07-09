package iLogger

import (
	"log"
	"os"
)

type Writer interface {
	Printf(string, ...interface{})
}

var DefaultLogger = log.New(os.Stdout, "generic-pkg", log.LstdFlags)
