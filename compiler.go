package main

import (
	"asritha.dev/compiler/pkg/scanner"
)

type Compiler interface {
	scanner.Scanner
	Compile() // take in io.Reader? or file
}
