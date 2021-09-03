package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const example = `
location /$location_var{
	alias $path;
}

запуск 
	nginx_configurator -s example.conf -v location_var:testHost,path:/root
вернет

location /testHost{
	alias /root;
}`

func main() {
	SampleFile := flag.String("s", "alias.conf", "файл шаблона")
	OutputFile := flag.String("o", "", "output path")
	Example := flag.Bool("example", false, "помощь по утилите")
	vars := flag.String("v", "", "переменные в формате ключ:значение,ключ:значение")
	flag.Parse()

	if *Example {
		fmt.Println(example)
		return
	}

	patternBytes, err := ioutil.ReadFile(*SampleFile)
	if err != nil {
		log.Fatal(err)
	}
	pattern := string(patternBytes)

	var writer io.Writer
	if *OutputFile == "" {
		writer = os.Stdout
	} else {
		if _, err := os.Stat(*OutputFile); err == nil {
			log.Fatal("file already exist")
		} else if os.IsNotExist(err) {
			writer, err = os.OpenFile(*OutputFile, os.O_CREATE|os.O_SYNC, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	for _, duo := range strings.Split(*vars, ",") {
		splitDuo := strings.Split(duo, ":")
		if len(splitDuo) != 2 {
			continue
		}

		key, value := splitDuo[0], splitDuo[1]

		pattern = strings.ReplaceAll(pattern,"$"+key,value)
	}

	fmt.Fprint(writer,string(pattern))
}
