package main

import (
	"os"
	"tweet-audit/internal/parser"
)

func main() {
	p, err := parser.NewContentParser("flagged.csv")
	if err != nil {
		panic(err)
	}
	path := os.Args[1]
	err = p.Parse(path)
	if err != nil {
		panic(err)
	}
}
