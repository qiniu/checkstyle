go-checkstyle
=============

checkstyle is a style check tool like java checkstyle. This tool inspired by [java checkstyle](https://github.com/checkstyle/checkstyle), [golint] (https://github.com/golang/lint). The style is according to [Go Code Review Comments](https://code.google.com/p/go-wiki/wiki/CodeReviewComments)

##to install, run
  go get github.com/qiniu/checkstyle/gocheckstyle

##to run
  gocheckstyle -config=.go_style dir1 dir2

##config is json file like following
```
{
    "file_line": 500,
    "func_line": 50,
    "params_num":4,
    "results_num":3,
    "formated": true,
    "pkg_name": true,
    "ignore":[
        "a/*",
        "b/*/c/*.go"
    ],
    "fatal":[
        "formated"
    ]
}

```

##add to makefile
```
check_go_style:
	bash -c "mkdir -p checkstyle; cd checkstyle && export GOPATH=`pwd`/checkstyle && go get github.com/qiniu/checkstyle/gocheckstyle"
	checkstyle/bin/gocheckstyle -config=.go_style dir1 dir2

```


Run checkstyle with one or more filenames or directories. The output of this tool is a list of suggestions. If you need to force obey the rule, place it in fatal.

## checkstyle's difference with other tools
Checkstyle differs from gofmt. Gofmt reformats Go source code, whereas checkstyle prints out coding style suggerstion.

Checkstyle differs from golint. Checkstyle only checks coding rule, not mistake; it is just subjective suggestion.
