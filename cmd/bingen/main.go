package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	typeNames      = flag.String("type", "", "comma-separated list of type names; must be set")
	includePrivate = flag.Bool("private", false, "include private fields")
	output         = flag.String("output", "", "output file name; default srcdir/<type>_gobin.go")
)

func generate(fname string) (err error) {
	fInfo, err := os.Stat(fname)
	if err != nil {
		return err
	}

	g := Generator{
		GoFile: fname,
		IsDir:  fInfo.IsDir(),
		Types:  strings.Split(*typeNames, ","),
	}
	if err := g.Run(); err != nil {
		return fmt.Errorf("Error generating code: %v", err)
	}
	return nil
}

func main1() {
	gofile := os.Getenv("GOFILE")
	flag.Parse()
	if len(*typeNames) == 0 || gofile == "" {
		flag.Usage()
		os.Exit(2)
	}
	// types := strings.Split(*typeNames, ",")
	if err := generate(gofile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	f, err := os.Create("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "args: %v\n", os.Args)
	fmt.Fprintf(f, "env: %v\n", os.Environ())
	f.WriteString("Hello, World!")

}
