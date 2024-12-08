package benlog

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
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
	fmt.Printf("/----Timer Started\n")
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			fmt.Fprintf(oldStdout, "| %s\n", scanner.Text())
		}
	}()
	start := time.Now()
	f()
	d := time.Since(start)
	w.Sync()
	w.Close()
	os.Stdout = oldStdout
	wg.Wait()
	fmt.Printf("\\----Took %v\n", d.Round(100*time.Nanosecond))
}
