package inputs

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
)

//go:embed puzzle_inputs/**/*
var puzzleInputs embed.FS

var (
	useExample = flag.Bool("use_example", false, "Use the example input")
)

func findDay() int {
	skip := 0
	for {
		_, f, _, _ := runtime.Caller(skip)
		cmdDir := filepath.Base(filepath.Dir(filepath.Dir(f)))
		if cmdDir == "cmd" {
			dir := filepath.Base(filepath.Dir(f))
			day, err := strconv.Atoi(dir)
			if err == nil {
				return day
			}
		}
		skip++
	}
}

func File() fs.File {
	day := findDay()

	if *useExample {
		f, err := puzzleInputs.Open(fmt.Sprintf("puzzle_inputs/%d/example.txt", day))
		if err != nil {
			benlog.Exitf("%v", err)
		}
		return f
	}
	f, err := puzzleInputs.Open(fmt.Sprintf("puzzle_inputs/%d/input.txt", day))
	if err != nil {
		benlog.Exitf("%v", err)
	}
	return f
}

func String() string {
	f := File()
	defer f.Close()
	var sb strings.Builder
	if _, err := io.Copy(&sb, f); err != nil {
		benlog.Exitf("%v", err)
	}
	return sb.String()
}

func Bytes() []byte {
	f := File()
	defer f.Close()
	var sb bytes.Buffer
	if _, err := io.Copy(&sb, f); err != nil {
		benlog.Exitf("%v", err)
	}
	return sb.Bytes()
}
