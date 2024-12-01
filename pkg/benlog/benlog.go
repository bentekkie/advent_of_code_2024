package benlog

import (
	"log"
	"os"
)

func Exitf(format string, args ...any) {
	log.Printf("%v", args...)
	os.Exit(1)
}
