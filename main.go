package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const example = `
location /%location{
	alias %path;
}

запуск 
	nginx_configurator -s=.\templates\alias.conf -o=.\out\alias.conf -location=test -path=root
вернет

location /test{
	alias /root;
}`

func main() {
	var Example bool
	var SampleFile string
	var OutputFile string
	var Vars = make(map[string]string, 0)
	for _, arg := range os.Args[1:] {
		arg = strings.TrimLeft(arg, "-")
		split := strings.Split(arg, "=")
		switch split[0] {
		case "example":
			Example = true
		case "s":
			if len(split) == 2 {
				SampleFile = split[1]
			}
		case "o":
			if len(split) == 2 {
				OutputFile = split[1]
			}
		default:
			if len(split) == 2 {
				Vars[split[0]] = split[1]
			}
		}
	}

	if Example {
		fmt.Println(example)
		return
	}

	patternBytes, err := ioutil.ReadFile(SampleFile)
	if err != nil {
		log.Fatal(err)
	}
	pattern := string(patternBytes)

	var writer io.Writer
	if OutputFile == "" {
		writer = os.Stdout
	} else {
		if _, err := os.Stat(OutputFile); err == nil {
			log.Fatal("file already exist")
		} else if os.IsNotExist(err) {
			writer, err = os.OpenFile(OutputFile, os.O_CREATE|os.O_SYNC, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for key, value := range Vars {

		pattern = strings.ReplaceAll(pattern, "%"+key, value)
	}

	fmt.Fprint(writer, string(pattern))
}
