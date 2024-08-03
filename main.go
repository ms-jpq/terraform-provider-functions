package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Panicln(err)
	}
	dir := filepath.Dir(self)
	if err := os.Chdir(dir); err != nil {
		log.Panicln(err)
	}
	fmt.Println(self)
}
