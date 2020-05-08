package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ninedraft/enum2go/pkg/generator"
	"github.com/ninedraft/enum2go/pkg/static"

	"github.com/spf13/pflag"
)

func main() {
	pflag.CommandLine.SetOutput(os.Stdout)

	var dir string
	pflag.StringVarP(&dir, "dir", "d", "", `package dir to parse`)
	var targetFile string
	pflag.StringVarP(&targetFile, "out", "o", "enums_generated.go", `file to generated result`)

	pflag.Usage = func() {
		printLines("enum2go",
			pad(2, static.Usage),
			"Flags:",
		)
		pflag.PrintDefaults()
	}

	pflag.Parse()

	var cfg = &generator.Config{
		Dir:        dir,
		TargetFile: targetFile,
	}
	if err := generator.Run(cfg); err != nil {
		log.Fatal(err)
	}
}

func pad(n int, text string) string {
	var pad = strings.Repeat(" ", n)
	var lines = strings.Split(text, "\n")
	var buf = &strings.Builder{}
	buf.Grow(len(text) + len(lines))
	for _, line := range lines {
		_, _ = buf.WriteString(pad)
		_, _ = buf.WriteString(line)
		_, _ = buf.WriteString("\n")
	}
	return buf.String()
}

func printLines(header string, lines ...string) {
	fmt.Printf("%s\n\n", header)
	for _, line := range lines {
		fmt.Println(line)
	}
}
