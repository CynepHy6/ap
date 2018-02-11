package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func unpackStruct(arr []commonStruct, tail, dir string) {
	if len(arr) > 0 {
		createDir(filepath.Join(outputName, separator))
		if singleSeparate && (dir == dirLang || dir == dirParam) {
			byteValue, _ := json.MarshalIndent(arr, "", "    ")
			value := string(byteValue)
			name := dir + eJSON
			writeFileString(name, value)
		} else {
			createDir(filepath.Join(outputName, dir))
			for _, c := range arr {
				value := c.Value
				if len(c.Columns) > 0 {
					value = c.Columns
				}
				if len(c.Trans) > 0 {
					value = c.Trans
				}
				name := c.Name + tail
				if len(c.Table) > 0 {
					name = c.Table
				}
				name = filepath.Join(dir, name)
				writeFileString(name, value)
			}
		}
	}
}

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
	if sufMode {
		unpackStruct(file.Contracts, eSIM, dirCon)
		unpackStruct(file.Menus, _menu+ePTL, dirMenu)
		unpackStruct(file.Blocks, _block+ePTL, dirBlock)
		unpackStruct(file.Pages, ePTL, dirPage)
		unpackStruct(file.Tables, _table+eJSON, dirTable)
		unpackStruct(file.Parameters, _param+eCSV, dirParam)
		unpackStruct(file.Languages, _lang+eJSON, dirLang)
	} else {
		unpackStruct(file.Contracts, eSIM, dirCon)
		unpackStruct(file.Menus, ePTL, dirMenu)
		unpackStruct(file.Blocks, ePTL, dirBlock)
		unpackStruct(file.Pages, ePTL, dirPage)
		unpackStruct(file.Tables, eJSON, dirTable)
		unpackStruct(file.Parameters, eCSV, dirParam)
		unpackStruct(file.Languages, eJSON, dirLang)
	}

	if len(file.Data) > 0 {
		createDir(filepath.Join(outputName, dirData))
		for _, c := range file.Data {
			name := c.Table + eJSON
			if sufMode {
				name = c.Table + _data + eJSON
			}
			name = filepath.Join(dirData, name)
			result, _ := json.MarshalIndent(c, "", "    ")
			writeFileString(name, string(result))
		}
	}
	writeConfig(bs)
}
