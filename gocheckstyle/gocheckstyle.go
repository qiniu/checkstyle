package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/qiniu/checkstyle"
)

var config = flag.String("config", "", "config json file")
var reporterOption = flag.String("reporter", "plain", "report output format, plain or xml")

var checker checkstyle.Checker
var reporter Reporter

type Ignore struct {
	Files []string `json:"ignore"`
}

var ignore Ignore

type Reporter interface {
	ReceiveProblems(checker checkstyle.Checker, file string, problems []checkstyle.Problem)
	Report()
}

type plainReporter struct {
	normalProblems []*checkstyle.Problem
	fatalProblems  []*checkstyle.Problem
}

func (_ *plainReporter) printProblems(ps []*checkstyle.Problem) {
	for _, p := range ps {
		log.Printf("%v: %s\n", p.Position, p.Description)
	}
}

func (p *plainReporter) Report() {
	if len(p.normalProblems) != 0 {
		log.Printf(" ========= There are %d normal problems ========= \n", len(p.normalProblems))
		p.printProblems(p.normalProblems)
	}

	if len(p.fatalProblems) != 0 {
		log.Printf(" ========= There are %d fatal problems ========= \n", len(p.fatalProblems))
		p.printProblems(p.fatalProblems)
		os.Exit(1)
	}
	if len(p.normalProblems) == 0 && len(p.fatalProblems) == 0 {
		log.Println(" ========= There are no problems ========= ")
	}
}

func (p *plainReporter) ReceiveProblems(checker checkstyle.Checker, file string, problems []checkstyle.Problem) {
	for i, problem := range problems {
		if checker.IsFatal(&problem) {
			p.fatalProblems = append(p.fatalProblems, &problems[i])
		} else {
			p.normalProblems = append(p.normalProblems, &problems[i])
		}
	}
}

type xmlReporter struct {
	problems map[string][]checkstyle.Problem
	hasFatal bool
}

func (x *xmlReporter) printProblems(ps []checkstyle.Problem) {
	format := "\t\t<error line=\"%d\" column=\"%d\" severity=\"%s\" message=\"%s\" source=\"checkstyle.%s\" />\n"
	for _, p := range ps {
		severity := "warning"
		if checker.IsFatal(&p) {
			severity = "error"
			x.hasFatal = true
		}
		log.Printf(format, p.Position.Line, p.Position.Column, severity, p.Description, p.Type)
	}
}

func (x *xmlReporter) Report() {
	log.SetFlags(0)
	log.Print(xml.Header)
	log.Println(`<checkstyle version="4.3">`)
	for k, v := range x.problems {
		log.Printf("\t<file name=\"%s\">\n", k)
		x.printProblems(v)
		log.Println("\t</file>")
	}
	log.Println("</checkstyle>")
	if x.hasFatal {
		os.Exit(1)
	}
}

func (x *xmlReporter) ReceiveProblems(checker checkstyle.Checker, file string, problems []checkstyle.Problem) {
	if len(problems) == 0 {
		return
	}
	x.problems[file] = problems
}

func main() {
	flag.Parse()

	files := flag.Args()

	if config == nil {
		log.Fatalln("No config")
	}
	if reporterOption == nil || *reporterOption != "xml" {
		reporter = &plainReporter{}
	} else {
		reporter = &xmlReporter{problems: map[string][]checkstyle.Problem{}}
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

	reporter.Report()
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

	reporter.ReceiveProblems(checker, fileName, ps)
}

func isIgnoreFile(fileName string) bool {
	for _, v := range ignore.Files {
		if ok, _ := filepath.Match(v, fileName); ok {
			return true
		}
	}
	return false
}

func isIgnoreDir(dir string) bool {
	for _, v := range ignore.Files {
		if ok, _ := filepath.Match(v, dir); ok {
			return true
		}
	}
	return false
}

func checkDir(dir string) {
	if isIgnoreDir(dir) {
		return
	}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() && isIgnoreDir(path) {
			return filepath.SkipDir
		}
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") && !isIgnoreFile(path) {
			checkFile(path)
		}
		return err
	})
}
