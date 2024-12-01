package main

import (
	_ "embed"
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
)

//go:embed main.go.tmpl
var mainTemplateStr string

var day = flag.Int("day", 0, "The day to generate")

func main() {
	flag.Parse()
	if *day == 0 {
		benlog.Exitf("Must specify a day")
	}
	_, f, _, _ := runtime.Caller(0)
	cmdDir := filepath.Dir(filepath.Dir(f))
	if *day == -1 {
		for i := 1; i <= 25; i++ {
			genDay(cmdDir, i, false)
		}
	} else {
		genDay(cmdDir, *day, true)
	}
}
func genDay(cmdDir string, d int, delete bool) {
	if delete {
		os.RemoveAll(filepath.Join(cmdDir, strconv.Itoa(d)))
	}
	os.Mkdir(filepath.Join(cmdDir, strconv.Itoa(d)), 0755)
	main := filepath.Join(cmdDir, strconv.Itoa(d), "main.go")
	if !delete {
		if _, err := os.Stat(main); err == nil {
			return
		}
	}
	os.WriteFile(filepath.Join(cmdDir, strconv.Itoa(d), "main.go"), []byte(mainTemplateStr), 0644)
}
