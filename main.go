package main

import (
	"asritha.dev/compiler/pkg/scanner"
	"fmt"
)

func main() {
	tInitInfos := []*scanner.TokenInitInfo{
		{
			3,
			"FOR",
			"for",
		},
	}

	s := scanner.NewScanner(tInitInfos)
	tokens, err := s.Scan()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(tokens)
	}
}
