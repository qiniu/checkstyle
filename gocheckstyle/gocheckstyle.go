package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/checkstyle"
)

var config = flag.String("config", "", "config json file")

var checker checkstyle.Checker

type Ignore struct {
	Files []string `json:"ignore"`
}

var ignore Ignore

var normalProblems []*checkstyle.Problem
var fatalProblems []*checkstyle.Problem

func main() {
	flag.Parse()

	files := flag.Args()

	if config == nil {
		log.Fatalln("No config")
	}
	conf, err := ioutil.ReadFile(*config)
	if err != nil {
		log.Fatalf("Open config %v fail %v\n", *config, err)
	}

	err = json.Unmarshal(conf, &ignore)
	if err != nil {
		log.Fatalf("Parse config %v fail \n", *config, err)
	}
	checker, err = checkstyle.New(conf)
	if err != nil {
		log.Fatalf("New checker fail %v\n", err)
	}

	for _, v := range files {
		if isDir(v) {
			checkDir(v)
		} else {
			checkFile(v)
		}
	}

	if len(normalProblems) != 0 {
		log.Printf(" ========= There are %d normal problems ========= \n", len(normalProblems))
		printProblems(normalProblems)
	}

	if len(fatalProblems) != 0 {
		log.Printf(" ========= There are %d fatal problems ========= \n", len(fatalProblems))
		printProblems(fatalProblems)
		os.Exit(1)
	}
}

func printProblems(ps []*checkstyle.Problem) {
	for _, p := range ps {
		log.Printf("%v: %s\n", p.Position, p.Description)
	}
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func checkFile(fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Read File Fail %v %v\n", fileName, err)
	}

	ps, err := checker.Check(fileName, file)
	if err != nil {
		log.Fatalf("Parse File Fail %v %v\n", fileName, err)
	}

	for _, p := range ps {
		if checker.IsFatal(&p) {
			fatalProblems = append(fatalProblems, &p)
		} else {
			normalProblems = append(normalProblems, &p)
		}
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
