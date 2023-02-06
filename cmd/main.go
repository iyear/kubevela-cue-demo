package main

import (
	"fmt"
	"go.uber.org/multierr"
	"kubecue"
	"os"
	"time"
)

const file = "test/struct.go"

// working dir is root of project
func main() {
	start := time.Now()
	if err := do(); err != nil {
		panic(err)
	}
	fmt.Printf("time: %s", time.Since(start))
}

func do() (rerr error) {
	g, err := kubecue.NewGenerator(file)
	if err != nil {
		return err
	}

	f, err := os.Create("test/struct.cue")
	defer multierr.AppendInvoke(&rerr, multierr.Close(f))

	return g.Generate(f)
}
