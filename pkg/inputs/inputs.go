package inputs

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"iter"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
	"github.com/bentekkie/advent_of_code_2024/pkg/flags"
)

//go:embed puzzle_inputs/**/*
var puzzleInputs embed.FS

var findDay = sync.OnceValue(func() int {
	skip := 0
	for {
		_, f, _, ok := runtime.Caller(skip)
		if !ok {
			benlog.Exitf("Could not find day")
		}
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
})

func File() fs.File {
	day := findDay()

	if flags.UseExample() {
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

func Lines() iter.Seq[string] {
	f := File()
	scanner := bufio.NewScanner(f)
	return func(yield func(string) bool) {
		defer f.Close()
		for scanner.Scan() {
			if !yield(scanner.Text()) {
				return
			}
		}
		if err := scanner.Err(); err != nil {
			benlog.Exitf("%v", err)
		}
	}
}
