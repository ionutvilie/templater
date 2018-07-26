package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
	yaml "gopkg.in/yaml.v2"
)

var (
	err         error
	wr          *os.File
	values      map[string]interface{}
	valueFile   = kingpin.Flag("valueFile", "A .yaml file containing data").Default("values.yaml").String()
	templateDir = kingpin.Flag("tmplDir", "The golang text/template folder").Default("templates").String()
	outDir      = kingpin.Flag("outDir", "Where to generate out files").Default("out").String()
	dryRun      = kingpin.Flag("dry-run", "Prints to stdout").Default("false").Bool()
)

// Config is a structure that holds the app config
type Config struct {
	DryRun      bool
	ValueFile   string
	Values      *map[string]interface{} // not required
	TemplateDir string
	OutDir      string
	Templates   []Template
}

// Template is a structure that holds template details
type Template struct {
	InFile  string // not required
	InDir   string // not required
	OutFile string
	OutDir  string
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

// writeTemplateToFile safeley renders input template
// if a template fails to execute, it recovers and continues.
func writeTemplateToFile(templateFile string, values map[string]interface{}, wr io.Writer) error {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			log.Errorf("Error parsing template: %v", templateFile)
		}
	}()
	t, _ := template.New("template").Delims("{{{", "}}}").ParseFiles(templateFile)
	_ = t.ExecuteTemplate(wr, "Template", values)
	if err != nil {
		recover()
	}
	return nil
}

// newfile returns a file and created the parent dir if IsNotExist
func newFile(path, file string) *os.File {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	checkErr(err)
	f, err := os.Create(filepath.Join(path, file))
	checkErr(err)

	f, err = os.OpenFile(filepath.Join(path, file),
		os.O_RDWR|os.O_APPEND, os.ModePerm) // For read access.
	checkErr(err)

	return f
}

func (c *Config) getTemplatesFiles() {
	// fmt.Printf("config: %#v\n", c)
	err := filepath.Walk(c.TemplateDir, func(path string, info os.FileInfo, err error) error {
		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		if ext == ".tmpl" {
			outFolder := filepath.Join(*outDir, strings.Replace(dir, *templateDir, "", 1))
			outFile := filepath.Join(strings.Replace(file, ext, "", 1))
			t := Template{
				InDir:   dir,
				InFile:  file,
				OutDir:  outFolder,
				OutFile: outFile,
			}
			c.Templates = append(c.Templates, t)
			// fmt.Printf("config: %#v\n", t)
		}
		return nil
	})
	checkErr(err)
	// fmt.Printf("config: %#v\n", c)
}

func (c *Config) execTemplates() {
	for _, template := range c.Templates {
		if *dryRun == true {
			wr = os.Stdout
		} else {
			wr = newFile(template.OutDir, template.OutFile)
		}
		log.Printf("Exec %v", filepath.Join(template.InDir, template.InFile))
		err := writeTemplateToFile(filepath.Join(template.InDir, template.InFile), *c.Values, wr)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	defer wr.Close()

}

func init() {
	kingpin.Version("0.0.1")
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()
}

func main() {

	valuesFile, err := ioutil.ReadFile(*valueFile)
	checkErr(err)

	err = yaml.Unmarshal(valuesFile, &values)
	checkErr(err)

	config := Config{
		DryRun:      *dryRun,
		ValueFile:   *valueFile,
		Values:      &values,
		TemplateDir: *templateDir,
		OutDir:      *outDir,
	}
	// fmt.Printf("config: %#v\n", config)
	config.getTemplatesFiles()
	// fmt.Printf("config: %#v\n", config)
	config.execTemplates()

}
