package main

import (
	"flag"
	"log"
	"os"
	"path"
)

var filedir = flag.String("path", ".", "pass the path you want to create file")

func main() {
	flag.Parse()
	if err := os.MkdirAll(*filedir, os.ModePerm); err != nil {
		log.Fatal("error creating dir")
	}
	file := path.Join(*filedir, "migration.ok")
	if _, err := os.Create(file); err != nil {
		log.Fatal("Error creating file: " + err.Error())
	}
}
