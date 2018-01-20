package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	condition = "ContractConditions(\"MainCondition\")"
	menu      = "default_menu"
	sim       = ".sim"
	ptl       = ".ptl"
	outName   = "out.json"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"."}
	}
	for _, path := range args {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		content := createExportJSON(files)
		file, err := os.Create(outName)
		if err != nil {
			return
		}
		defer file.Close()
		file.WriteString(content)
	}
}

func createExportJSON(files []os.FileInfo) string {
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
	for _, file := range files {
		switch ext := filepath.Ext(file.Name()); ext {
		case ptl:
			pages = append(pages, convert(file.Name(), ptl))
		case sim:
			contracts = append(contracts, convert(file.Name(), sim))
		}
	}
	out["pages"] = pages
	out["contracts"] = contracts
	result, _ := json.Marshal(out)
	return string(result)
}

func convert(filename string, ext string) (result map[string]string) {
	result = make(map[string]string)
	result["Name"] = filename[:len(filename)-len(ext)]
	result["Conditions"] = condition
	switch ext {
	case ptl:
		result["Menu"] = menu
	}
	result["Value"] = file2JSON(filename)
	return
}

func file2JSON(filename string) (str string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	byteStr, _ := json.Marshal(string(bs))
	str = string(byteStr)
	return
}
