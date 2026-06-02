package main

import (
	"log"
	"os"
	"tweet-audit/internal/parser"
)

func main() {
	p, err := parser.NewContentParser("flagged.csv")
	if err != nil {
		panic(err)
	}
	if len(os.Args) < 2 {
		log.Fatal("missing required argument")
	}
	path := os.Args[1]
	err = p.Parse(path)
	if err != nil {
		panic(err)
	}
}
