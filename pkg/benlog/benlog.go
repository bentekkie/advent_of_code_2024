package benlog

import (
	"fmt"
	"log"
	"os"

	"github.com/bentekkie/advent_of_code_2024/pkg/flags"
)

func Exitf(format string, args ...any) {
	log.Printf("%v", args...)
	os.Exit(1)
}

func ExamplePrintf(format string, args ...any) {
	if flags.UseExample() {
		fmt.Printf(format, args...)
	}
}
