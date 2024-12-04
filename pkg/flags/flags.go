package flags

import "flag"

var (
	useExample = flag.Bool("use_example", false, "Use the example input")
)

func UseExample() bool {
	return *useExample
}
