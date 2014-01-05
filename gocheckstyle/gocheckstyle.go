package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/longbai/checkstyle"
)

var gConfig []byte
var config = flag.String("config", "", "config json file")

func main() {
	flag.Parse()

	files := flag.Args()
	for _, v := range files {
		if isDir(v) {
			checkDir(v)
		} else {
			checkFile(v)
		}
	}
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func checkFile(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Printf("Read File Fail %v %v", fileName, err)
	}

	checker, err := checkstyle.New(gConfig)
	if err != nil {
		panic("New Checker Fail")
	}

	ps, err := checker.Check(fileName, file)
	if err != nil {
		log.Printf("Parse File Fail %v %v", fileName, err)
	}

	for _, p := range ps {
		fmt.Printf("%s:%v: %s\n", fileName, p.Position, p.Description)
	}
}

func checkDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			checkFile(path)
		}
		return err
	})
}
