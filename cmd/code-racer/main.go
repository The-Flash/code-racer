package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/The-Flash/code-racer/internal/manifest"
)

var (
	manifestPtr = flag.String("f", "", "Path to manifest file")
)

func main() {
	flag.Parse()
	manifestFile, err := os.ReadFile(*manifestPtr)
	if err != nil {
		log.Fatal("could not read manifest file", err)
	}

	m := new(manifest.Manifest)
	m.Load(manifestFile)

	fmt.Println(m.Runtimes)

}
