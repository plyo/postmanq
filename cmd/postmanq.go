package main

import (
	"github.com/AdOnWeb/postmanq"
	"flag"
)


func main() {
	var file string
	flag.StringVar(&file, "f", postmanq.EXAMPLE_CONFIG_YAML, "configuration yaml file")
	flag.Parse()

	app := postmanq.NewApplication()
	if app.IsValidConfigFilename(file) {
		app.ConfigFilename = file
		app.Run()
	} else {
		flag.PrintDefaults()
	}
}
