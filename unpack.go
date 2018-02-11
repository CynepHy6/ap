package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	graphFile map[string][]string
)

func init() {
	graphFile = map[string][]string{}
}

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
				name := c.Name
				if len(c.Table) > 0 {
					name = c.Table
				}
				var mName string
				switch dir {
				case dirPage:
					mName = name + _page
					graphFile[mName] = append(graphFile[mName], getRelations(value, dir)...)
				case dirCon:
					mName = name + _contr
					graphFile[mName] = append(graphFile[mName], getRelations(value, dir)...)
				case dirMenu:
					mName = name + _menu
					graphFile[mName] = append(graphFile[mName], getRelations(value, dir)...)
				}
				nameTail := name + tail
				nameTail = filepath.Join(dir, nameTail)
				writeFileString(nameTail, value)
			}
		}
	}
	storeGraphviz(graphFile)
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

func getRelations(s, dir string) (entries []string) {
	contr2Contr := regexp.MustCompile(`[^(Join|info|warning|error|LangRes|FindEcosystem|CallContract|ContractAccess|ContractConditions|EvalCondition|ValidateCondition|AddressToId|Contains|Float|HasPrefix|HexToBytes|Int|Len|PubToID|IdToAddress|Money|Replace|Size|Sha256|Sprintf|Str|Substr|UpdateLang|SysParamString|SysParamInt|UpdateSysParam|EcosysParam|DBFind|DBInsert|DBInsertReport|DBUpdate|DBUpdateExt|DBRow|DBIntExt|DBStringExt)]\s*\(@?.*?\)`)
	page2Contr := regexp.MustCompile(`\(.*?Contract:\s*(@?\w+)`)
	page2Page := regexp.MustCompile(`\(.*?Page:\s*(\w+)`)
	contr2Table := regexp.MustCompile(`(?:DBFind|DBInsert|DBUpdate|DBUpdateExt|DBRow)\s*\(\s*"([\w]+?)"`)
	page2Table := regexp.MustCompile(`(?:DBFind\s*\(\s*Name:\s*)(.*?)[,\s]`)
	page2TableShort := regexp.MustCompile(`(?:DBFind\s*\(\s*)([^:]*?)[,\s]`)
	var parts, parts2, parts3 []string
	switch dir {
	case dirCon:
		parts = contr2Contr.FindStringSubmatch(s)
		addTypeString(parts, _contr)
		parts2 = contr2Table.FindStringSubmatch(s)
		addTypeString(parts2, _table)
	case dirPage:
		parts = page2Contr.FindStringSubmatch(s)
		addTypeString(parts, _contr)
		parts2 = page2Table.FindStringSubmatch(s)
		addTypeString(parts2, _table)
		parts3 = page2TableShort.FindStringSubmatch(s)
		addTypeString(parts3, _table)
	case dirMenu:
		parts = page2Page.FindStringSubmatch(s)
		addTypeString(parts, _page)
	}

	if len(parts) > 0 {
		entries = append(entries, parts[1:]...)
	}
	if len(parts2) > 0 {
		entries = append(entries, parts2[1:]...)
	}
	if len(parts3) > 0 {
		entries = append(entries, parts3[1:]...)
	}
	fmt.Println(entries)
	return
}

func storeGraphviz(m map[string][]string) {
	resArr := []string{"digraph G{"}
	config := "rankdir=RL;"
	resArr = append(resArr, config)
	for name, rels := range m {
		var g string
		if len(rels) > 0 {
			g = fmt.Sprintf("%s -> %s", name, strings.Join(rels, " -> "))
		} else {
			g = fmt.Sprintf("%s", name)
		}
		resArr = append(resArr, g)
	}
	resArr = append(resArr, "}")
	out := strings.Join(resArr, "\n")

	writeFileString("graphviz.dot", out)
}
func addTypeString(arr []string, t string) []string {
	for i := range arr {
		arr[i] = fmt.Sprintf("%s%s", arr[i], t)
	}
	return arr
}
