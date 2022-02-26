package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/shuntagami/dojo1/kadai1/shuntagami/converter"
)

func main() {
	flag.Parse()
	from := "." + strings.ToLower(flag.Arg(0))
	to := "." + strings.ToLower(flag.Arg(1))
	targetDirName := flag.Arg(2)

	if err := converter.Initialize(from, to, os.Getenv("PROJECT_ROOT_DIR")); err != nil {
		log.Fatal(err)
	}

	if err := converter.Client.Convert(targetDirName); err != nil {
		log.Fatal(err)
	}
}
