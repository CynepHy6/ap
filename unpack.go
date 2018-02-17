package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func unpackJSON(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := importFile{}
	if err := json.Unmarshal(bs, &file); err != nil {
		fmt.Println("unmarshal file error:", err)
		return
	}

	unpackStruct(file.Contracts, eSIM, dirCon)
	unpackStruct(file.Menus, ePTL, dirMenu)
	unpackStruct(file.Blocks, ePTL, dirBlock)
	unpackStruct(file.Pages, ePTL, dirPage)
	unpackStruct(file.Tables, eJSON, dirTable)
	unpackStruct(file.Parameters, eCSV, dirParam)
	unpackStruct(file.Languages, eJSON, dirLang)

	if len(file.Data) > 0 {
		createDir(filepath.Join(outputName, dirData))
		for _, c := range file.Data {
			name := c.Table + eJSON
			name = filepath.Join(dirData, name)
			result, _ := json.MarshalIndent(c, "", "    ")
			writeFileString(name, string(result))
		}
	}
	writeConfig(bs)
	if abs, err := filepath.Abs(outputName); err == nil {
		abspath := filepath.Join(abs, structFileName)
		createGraph(abspath)
	}
}

func unpackStruct(arr []commonStruct, tail, dir string) {
	if len(arr) > 0 {
		createDir(filepath.Join(outputName, dir))
		for _, c := range arr {
			value := c.Value
			if len(c.Columns) > 0 {
				value = c.Columns
			}
			if len(c.Trans) > 0 {
				value = c.Trans
			}
			name := c.Name
			if len(c.Table) > 0 {
				name = c.Table
			}
			nameTail := name + tail
			nameTail = filepath.Join(dir, nameTail)
			writeFileString(nameTail, value)
		}
	}
}
