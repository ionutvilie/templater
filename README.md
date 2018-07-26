# templater

[![Go Report Card](https://goreportcard.com/badge/github.com/ionutvilie/templater)](https://goreportcard.com/report/github.com/ionutvilie/templater)

A simple golang cli app that accepts a yaml file as input and applies it over a golang templates folder

## usage 

```bash 
./templater -h
usage: main [<flags>]

Flags:
  -h, --help                     Show context-sensitive help (also try --help-long and --help-man).
      --valueFile="values.yaml"  A .yaml file containing data
      --tmplDir="templates"      The golang text/template folder
      --outDir="out"             Where to generate out files
      --dry-run                  Prints to stdout
      --version                  Show application version.
```

## how it works

```bash
# values.yaml --- >  /templates/phone-book.json.tmpl -- ( removes .tmpl ) --> /out/phone-book.json

# example template
cat templates/phone-book.json.tmpl
{{{ define "Title"}}}Humans{{{end}}}
{{{ define "Template" }}}{
    "title": "{{{ template "Title" }}}",
    "data": [ {{{ range $key, $human :=  .Humans }}}{{{if $key}}},{{{end}}}
    {   "firstName" : "{{{ $human.FirstName }}}",
        "lastName" : "{{{ $human.LastName }}}",
        "occupation" : "{{{ $human.Occupation }}}",
        "phone" : "{{{ $human.Phone }}}"
    }{{{ end }}}
    ]
}
{{{ end }}}

# example of output file
cat out/phone-book.json
{
    "title": "Humans",
    "data": [
    {   "firstName" : "James",
        "lastName" : "Darakjy",
        "occupation" : "Foo",
        "phone" : "504-621-8927"
    },
    {   "firstName" : "Mitsue",
        "lastName" : "Dilliard",
        "occupation" : "Bar",
        "phone" : "513-570-18939"
    }
    ]
}

```
