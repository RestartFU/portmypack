package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/restartfu/portmypack/portmypack"
	"github.com/restartfu/portmypack/portmypack/bedrock"
	"github.com/restartfu/portmypack/portmypack/java"
)

func main() {
	input := flag.String("i", "", "input Java Edition pack: a .zip file or a folder")
	output := flag.String("o", ".", "output directory; with -mcpack, the .mcpack file to write or the directory to write <pack>.mcpack in")
	mcpack := flag.Bool("mcpack", false, "write a single .mcpack file to -o instead of extracting the pack")
	minEngine := flag.String("min-engine", versionString(bedrock.DefaultMinEngineVersion),
		"oldest Bedrock version allowed to load the pack, as major.minor.patch")
	flag.Parse()
	if len(*input) == 0 || len(*output) == 0 {
		flag.Usage()
		return
	}
	engine, err := parseVersion(*minEngine)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid -min-engine:", err)
		os.Exit(2)
	}

	javapack, err := java.NewResourcePack(*input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	opt := portmypack.WithMinEngineVersion(engine)
	if *mcpack {
		dest := *output
		if info, statErr := os.Stat(dest); statErr == nil && info.IsDir() {
			dest = filepath.Join(dest, javapack.Name+".mcpack")
		}
		err = portmypack.PortJavaEditionPack(javapack, dest, opt)
	} else {
		err = portmypack.PortJavaEditionPackAndExtract(javapack, *output, opt)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// parseVersion parses a major.minor.patch version; missing parts are zero.
func parseVersion(s string) ([3]int, error) {
	var v [3]int
	parts := strings.Split(s, ".")
	if len(parts) == 0 || len(parts) > 3 {
		return v, fmt.Errorf("%q is not major.minor.patch", s)
	}
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil || n < 0 {
			return v, fmt.Errorf("%q is not major.minor.patch", s)
		}
		v[i] = n
	}
	return v, nil
}

// versionString formats a version as major.minor.patch.
func versionString(v [3]int) string {
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}
