package main

import (
	"flag"

	"github.com/restartfu/portmypack/portmypack"
	"github.com/restartfu/portmypack/portmypack/java"
)

func main() {
	input := flag.String("i", "", "input path")
	output := flag.String("o", ".", "output path")
	flag.Parse()
	if len(*input) == 0 || len(*output) == 0 {
		flag.Usage()
		return
	}

	javapack, err := java.NewResourcePack(*input)
	if err != nil {
		panic(err)
	}
	portmypack.PortJavaEditionPack(javapack, *output)
}
