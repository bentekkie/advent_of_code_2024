package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/bentekkie/advent_of_code_2024/pkg/benlog"
)

//go:embed session.txt
var sessionCookie string

var day = flag.Int("day", -1, "The day to generate")

func main() {
	flag.Parse()
	if *day == 0 {
		benlog.Exitf("Must specify a day")
	}
	_, thisFile, _, _ := runtime.Caller(0)
	// ../../../pkg
	pkgDir := filepath.Join(filepath.Dir(filepath.Dir(filepath.Dir(thisFile))), "pkg")
	loc, _ := time.LoadLocation("America/Toronto")
	if *day == -1 {
		now := time.Now().In(loc)
		for i := 1; i <= now.Day(); i++ {
			generateDay(pkgDir, i)
		}
	} else {
		generateDay(pkgDir, *day)
	}
}

func generateDay(pkgDir string, day int) {
	puzzleInputsDir := filepath.Join(pkgDir, "inputs", "puzzle_inputs", strconv.Itoa(day))
	os.MkdirAll(puzzleInputsDir, 0755)
	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/2024/day/%d/input", day), nil)
	if err != nil {
		benlog.Exitf("%v", err)
	}
	req.AddCookie(
		&http.Cookie{
			Name:  "session",
			Value: sessionCookie,
		},
	)
	req.Header.Set("User-Agent", "input fetching code for bentekkie@gmail.com")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		benlog.Exitf("%v", err)
	}
	defer resp.Body.Close()
	f, err := os.Create(filepath.Join(puzzleInputsDir, "input.txt"))
	if err != nil {
		benlog.Exitf("%v", err)
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		benlog.Exitf("%v", err)
	}
}
