package main

import (
	"flag"
	"fmt"
	"go.uber.org/multierr"
	"kubecue"
	"os"
	"path/filepath"
	"time"
)

var f = flag.String("file", "", "file to convert")

func init() {
	flag.Parse()
}

// working dir is root of project
func main() {
	start := time.Now()
	if err := do(*f); err != nil {
		panic(err)
	}
	fmt.Printf("time: %s", time.Since(start))
}

func do(file string) (rerr error) {
	g, err := kubecue.NewGenerator(file)
	if err != nil {
		return err
	}

	g.RegisterAny(
		"*k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.Unstructured",
	)

	f, err := os.Create(changeExt(file, ".cue"))
	defer multierr.AppendInvoke(&rerr, multierr.Close(f))

	return g.Generate(f)
}

func changeExt(file, ext string) string {
	old := filepath.Ext(file)
	return file[:len(file)-len(old)] + ext
}
