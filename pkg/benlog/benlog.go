package benlog

import (
	"fmt"
	"log"
	"os"
	"time"

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

func Timed(f func()) {
	start := time.Now()
	f()
	fmt.Printf("took %v\n", time.Since(start))
}
