go-checkstyle
=============

checkstyle is a style check tool like java checkstyle. This tool inspired by [java checkstyle](https://github.com/checkstyle/checkstyle), [golint] (https://github.com/golang/lint)

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
    "ignore":[
        "a/*",
        "b/*/c/*.go"
    ],
    "fatal":[
        "formated"
    ]
}

```

Run checkstyle with one or more filenames or directories. The output of this tool is a list of suggestions. If you need to force obey the rule, place it in fatal.

## checkstyle's difference with other tools
Checkstyle differs from gofmt. Gofmt reformats Go source code, whereas checkstyle prints out coding style suggerstion.

Checkstyle differs from golint. Checkstyle only checks coding conventions, not mistake; it is just subjective suggestion. In the future, I will integrate some golint rule. 




