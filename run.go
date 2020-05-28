package main

import (
	"go-wordfilter/wordfilter"
)

func main() {
	wordfilter.LoadWords()
	wordfilter.NewInterfaceHttp().Run()
}
