package main

import (
	"flag"
	"log"

	"github.com/ninedraft/enum2go/pkg/generator"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", "", `package dir to parse`)
	var targetFile string
	flag.StringVar(&targetFile, "out", "enums_generated.go", `file to generated result`)
	flag.Parse()

	var cfg = &generator.Config{
		Dir:             dir,
		TargetFile:      targetFile,
		TypePlaceholder: "Θ",
	}
	if err := generator.Run(cfg); err != nil {
		log.Fatal(err)
	}
}
