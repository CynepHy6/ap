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
		createExportJSON(files)
	}
}

func createExportJSON(files []os.FileInfo) {

	contracts := []map[string]string{}
	pages := []map[string]string{}
	out := make(map[string][]map[string]string)

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
	fmt.Println(string(result))
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
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("%s not found", filename)
		return
	}
	defer file.Close()

	// get stat
	stat, err := file.Stat()
	if err != nil {
		fmt.Printf("%s not stats", filename)
		return
	}
	// read
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Printf("%s not read", filename)
		return
	}

	byteStr, _ := json.Marshal(string(bs))
	str = string(byteStr)
	return
}
