package main

import (
	"flag"
	"fmt"
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

	fmt.Println("parsing file: ", *f)

	var err error
	switch flag.Arg(0) {
	case "provider":
		err = providerGen(*f)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("time: %s", time.Since(start))
}

func changeExt(file, ext string) string {
	old := filepath.Ext(file)
	return file[:len(file)-len(old)] + ext
}
