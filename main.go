// Copyright 2018 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// file2byteslice is a dead simple tool to embed a file to Go.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	inputFilename  = flag.String("input", "", "input file name")
	outputFilename = flag.String("output", "", "output file name")
	packageName    = flag.String("package", "main", "package name")
	varName        = flag.String("var", "_", "variable name")
)

func write(w io.Writer, r io.Reader) error {
	fmt.Fprintln(w, "// DO NOT EDIT - This file is generated by file2byteslice")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "package %s\n", *packageName)
	fmt.Fprintf(w, "var %s = []byte{\n", *varName)

	i := 0
	for {
		buf := make([]byte, 4096)
		n, err := r.Read(buf)

		for _, b := range buf[:n] {
			fmt.Fprintf(w, "0x%02x,", b)
			if (i+1)%16 == 0 {
				fmt.Fprintln(w)
			}
			i++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	fmt.Fprintln(w, "\n}")
	return nil
}

func run() error {
	var out io.Writer
	if *outputFilename != "" {
		f, err := os.Create(*outputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}

	var in io.Reader
	if *inputFilename != "" {
		f, err := os.Open(*inputFilename)
		if err != nil {
			return err
		}
		defer f.Close()
		in = f
	} else {
		in = os.Stdin
	}

	if err := write(out, in); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		panic(err)
	}
}
