go-checkstyle
=============
[![Build Status](https://api.travis-ci.org/qiniu/checkstyle.png?branch=master)](https://travis-ci.org/qiniu/checkstyle)

checkstyle is a style check tool like java checkstyle. This tool inspired by [java checkstyle](https://github.com/checkstyle/checkstyle), [golint](https://github.com/golang/lint). The style refered to some points in [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments).

# Install
  go get github.com/qiniu/checkstyle/gocheckstyle

# Run
  gocheckstyle -config=.go_style dir1 dir2

# Config 
config is json file like the following:
```
{
    "file_line": 500,
    "func_line": 50,
    "params_num":4,
    "results_num":3,
    "formated": true,
    "pkg_name": true,
    "camel_name":true,
    "ignore":[
        "a/*",
        "b/*/c/*.go"
    ],
    "fatal":[
        "formated"
    ]
}

```

# Add to makefile
```
check_go_style:
	bash -c "mkdir -p checkstyle; cd checkstyle && export GOPATH=`pwd` && go get github.com/qiniu/checkstyle/gocheckstyle"
	checkstyle/bin/gocheckstyle -config=.go_style dir1 dir2

```

# Integrate with jenkins checkstyle plugin
excute in shell
```
    mkdir -p checkstyle; cd checkstyle && export GOPATH=`pwd` && go get github.com/qiniu/checkstyle/gocheckstyle"
    checkstyle/bin/gocheckstyle -reporter=xml -config=.go_style dir1 dir2 2>gostyle.xml
```
then add postbuild checkstyle file gostyle.xml

Run checkstyle with one or more filenames or directories. The output of this tool is a list of suggestions. If you need to force obey the rule, place it in fatal.

# Checkstyle's difference with other tools
Checkstyle differs from gofmt. Gofmt reformats Go source code, whereas checkstyle prints out coding style suggestion.

Checkstyle differs from golint. Checkstyle check file line/function line/param number, could be configed by user.
