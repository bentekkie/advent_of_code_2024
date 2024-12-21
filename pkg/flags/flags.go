package flags

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

var (
	useExample = flag.Bool("use_example", false, "Use the example input")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func UseExample() bool {
	return *useExample
}

func CPUProfile() func() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		return pprof.StopCPUProfile
	}
	return func() {}
}
