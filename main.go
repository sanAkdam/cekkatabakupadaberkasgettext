package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"sync"
)

var (
	wg sync.WaitGroup
)

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		wg.Add(1)
		go performGettextCheck(arg)
	}
	wg.Wait()
	fmt.Println("ok")
}

func performGettextCheck(filename string) {
	defer wg.Done()
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	translations := parseTranslation(string(file))
	ctx := NewContext(filename, translations)
	for k, v := range ctx.Run() {
		fmt.Println(k + " => " + v)
	}
}

func parseTranslation(file string) (output []string) {
	reg := regexp.MustCompile("msgstr \"(.+)\"")
	for _, t := range reg.FindAllStringSubmatch(string(file), -1) {
		output = append(output, t[1])
	}
	return
}
