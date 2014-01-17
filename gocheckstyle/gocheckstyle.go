package main

import (
	"encoding/json"
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

var config = flag.String("config", "", "config json file")

var checker checkstyle.Checker

type Ignore struct {
	Files []string `json:"ignore"`
}

var ignore Ignore

func main() {
	flag.Parse()

	files := flag.Args()

	if config == nil {
		log.Println("No Config")
		return
	}
	conf, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Printf("Open Config %v Fail %v\n", *config, err)
		return
	}

	err = json.Unmarshal(conf, &ignore)
	if err != nil {
		log.Printf("Parse Config %v Fail \n", *config, err)
		return
	}
	checker, err = checkstyle.New(conf)
	if err != nil {
		fmt.Errorf("New Checker Fail %s\n", err.Error())
		return
	}

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
		log.Printf("Read File Fail %v %v\n", fileName, err)
	}

	ps, err := checker.Check(fileName, file)
	if err != nil {
		log.Printf("Parse File Fail %v %v\n", fileName, err)
	}

	for _, p := range ps {
		fmt.Printf("%v: %s\n", p.Position, p.Description)
	}
}

func isIgnore(fileName string) bool {
	for _, v := range ignore.Files {
		if v == fileName {
			return true
		}
	}
	return false
}

func checkDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") && !isIgnore(path) {
			checkFile(path)
		}
		return err
	})
}
