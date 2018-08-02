package main

import (
	"io"
	"log"
	"os"

	"gopl.io/ch10/file/extract"
	_ "gopl.io/ch10/file/filetar"
	_ "gopl.io/ch10/file/filezip"
)

func main() {
	handleFile := func(reader io.Reader) {
		if _, err := io.Copy(os.Stdout, reader); err != nil {
			log.Fatal(err)
		}
	}

	if err := extract.Read(os.Args[1], handleFile); err != nil {
		log.Fatal(err)
	}
}
