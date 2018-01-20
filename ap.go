package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

const (
	sim = ".sim"
	ptl = ".ptl"
)

var (
	condition string
	menu      string
	outName   string
	prefix    string
	path      string
)

func init() {
	flag.StringVar(&condition, "c", "ContractConditions(\"MainCondition\")", "shortcut for for --conditions")
	flag.StringVar(&condition, "-conditions", "ContractConditions(\"MainCondition\")", "conditions")
	flag.StringVar(&menu, "m", "default_menu", "shortcut for --menu")
	flag.StringVar(&menu, "-menu", "default_menu", "menu")
	flag.StringVar(&outName, "o", "out", "shortcut for --output")
	flag.StringVar(&outName, "-output", "out", "output filename for JSON")
	flag.StringVar(&prefix, "p", "", "shortcut for --prefix")
	flag.StringVar(&prefix, "-prefix", "", "prefix for pages and contracts")
	flag.StringVar(&path, "i", ".", "shortcut for --input")
	flag.StringVar(&path, "-input", ".", "path for input files")
}

func main() {
	flag.Parse()
	if prefix != "" {
		prefix = prefix + "_"
	}
	content := createJSON(path)
	outFile, err := os.Create(outName + ".json")
	if err != nil {
		return
	}
	defer outFile.Close()
	outFile.WriteString(content)
}

func createJSON(path string) string {
	emptyEntry := []map[string]string{}
	contracts := emptyEntry
	pages := emptyEntry
	out := make(map[string][]map[string]string)
	out["menus"] = emptyEntry
	out["parameters"] = emptyEntry
	out["languages"] = emptyEntry
	out["tables"] = emptyEntry
	out["data"] = emptyEntry
	out["blocks"] = emptyEntry
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	for _, file := range files {
		fileName := file.Name()
		switch ext := filepath.Ext(fileName); ext {
		case ptl:
			pages = append(pages, convert(path, fileName, ptl))
		case sim:
			contracts = append(contracts, convert(path, fileName, sim))
		}
	}
	out["pages"] = pages
	out["contracts"] = contracts
	result, _ := json.Marshal(out)
	return string(result)
}

func convert(path, filename, ext string) (result map[string]string) {
	result = make(map[string]string)
	name := filename[:len(filename)-len(ext)]
	result["Name"] = prefix + name
	result["Conditions"] = condition
	result["Value"] = file2str(path, filename)
	switch ext {
	case sim:
		if prefix != "" {
			re := regexp.MustCompile("contract\\s+" + name)
			result["Value"] = re.ReplaceAllString(result["Value"], "contract "+result["Name"])
		}
	case ptl:
		result["Menu"] = menu
	}
	return
}

func file2str(path, filename string) (str string) {
	bs, err := ioutil.ReadFile(path + filename)
	if err != nil {
		return
	}
	str = string(bs)
	return
}
