package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tmc/dot"
)

const (
	structFileName = "struct.dot"
	pageColor      = "green"
	contrColor     = "red"
	menuColor      = "blue"
)

var (
	graphMap     = map[string][]string{}
	graphDot     = dot.NewGraph("G")
	contr2Contr  = regexp.MustCompile("[^(Join|info|warning|error|LangRes|FindEcosystem|CallContract|ContractAccess|ContractConditions|EvalCondition|ValidateCondition|AddressToId|Contains|Float|HasPrefix|HexToBytes|Int|Len|PubToID|IdToAddress|Money|Replace|Size|Sha256|Sprintf|Str|Substr|UpdateLang|SysParamString|SysParamInt|UpdateSysParam|EcosysParam|DBFind|DBInsert|DBInsertReport|DBUpdate|DBUpdateExt|DBRow|DBIntExt|DBStringExt)]\\s*\\(@?.*?\\)")
	page2Contr   = regexp.MustCompile("\\(.*?Contract:\\s*(@?\\w+)")
	page2Page    = regexp.MustCompile("\\(.*?Page:\\s*(\\w+)")
	contr2Table  = regexp.MustCompile("(?:DBFind|DBInsert|DBUpdate|DBUpdateExt|DBRow)\\s*\\(\\s*[\"\\`]([\\w]+?)\"")
	page2Table   = regexp.MustCompile("DBFind\\s*\\(\\s*Name:\\s*(.*?)[,\\s]|DBFind\\s*\\(\\s*([^:]*?)[\\),\\s]")
	includeBlock = regexp.MustCompile("Include\\s*\\(\\s*Name:\\s*(.*?)[,\\s]|Include\\s*\\(\\s*([^:]*?)[\\),\\s]")
)

func unpackJSON(filename string) {
	graphDot.SetType(dot.DIGRAPH)
	graphDot.Set("rankdir", "LR")
	graphDot.Set("fontsize", "20.0")
	labelGraph := fmt.Sprintf("%s %s", strings.Trim(outputName, separator), time.Now().Format(time.RFC850))
	graphDot.Set("label", labelGraph)

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
		unpackStruct(file.Contracts, eSIM, dirCon, graphDot)
		unpackStruct(file.Menus, _menu+ePTL, dirMenu, graphDot)
		unpackStruct(file.Blocks, _block+ePTL, dirBlock, graphDot)
		unpackStruct(file.Pages, ePTL, dirPage, graphDot)
		unpackStruct(file.Tables, _table+eJSON, dirTable, graphDot)
		unpackStruct(file.Parameters, _param+eCSV, dirParam, graphDot)
		unpackStruct(file.Languages, _lang+eJSON, dirLang, graphDot)
	} else {
		unpackStruct(file.Contracts, eSIM, dirCon, graphDot)
		unpackStruct(file.Menus, ePTL, dirMenu, graphDot)
		unpackStruct(file.Blocks, ePTL, dirBlock, graphDot)
		unpackStruct(file.Pages, ePTL, dirPage, graphDot)
		unpackStruct(file.Tables, eJSON, dirTable, graphDot)
		unpackStruct(file.Parameters, eCSV, dirParam, graphDot)
		unpackStruct(file.Languages, eJSON, dirLang, graphDot)
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
	writeFileString(structFileName, graphDot.String())
}

func unpackStruct(arr []commonStruct, tail, dir string, graph *dot.Graph) {
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
				switch dir { // parse graph
				case dirPage:
					fallthrough
				case dirCon:
					fallthrough
				case dirTable:
					fallthrough
				case dirBlock:
					fallthrough
				case dirMenu:
					node := dot.NewNode(getNodeName(name, dir))
					if dir == dirPage || dir == dirBlock {
						node.Set("fontcolor", pageColor)
					}
					if dir == dirCon {
						node.Set("fontcolor", contrColor)
					}
					if dir == dirMenu {
						node.Set("fontcolor", menuColor)
					}
					node.Set("group", dir)
					if dir != dirTable {
						addEdges(node, value, dir)
					}
					graph.AddNode(node)
				}
				nameTail := name + tail
				nameTail = filepath.Join(dir, nameTail)
				writeFileString(nameTail, value)
			}
		}
	}

}

func addEdges(parentNode *dot.Node, s, dir string) {
	switch dir {
	case dirCon:
		addNode(parentNode, contr2Contr, s, dir, "")
		addNode(parentNode, contr2Table, s, dirTable, "")
	case dirPage:
		addNode(parentNode, page2Contr, s, dirCon, "")
		addNode(parentNode, page2Table, s, dirTable, "")
		addNode(parentNode, page2Page, s, dir, "")
		addNode(parentNode, includeBlock, s, dirBlock, "Include")
	case dirBlock:
		addNode(parentNode, page2Contr, s, dirCon, "")
		addNode(parentNode, page2Table, s, dirTable, "")
		addNode(parentNode, page2Page, s, dir, "")
	case dirMenu:
		addNode(parentNode, page2Page, s, dirPage, "")
		// fmt.Println(graphDot)
	}
}

func addNode(parentNode *dot.Node, pat *regexp.Regexp, str, dir, label string) {
	s := strings.Replace(str, "`", `"`, -1)
	arr := pat.FindAllStringSubmatch(s, -1)
	for _, match := range arr {
		for i := range match {
			if i > 0 {
				if match[i] != "" {
					name := getNodeName(match[i], dir)
					if !stringInSlice(graphMap[parentNode.Name()], name) { // check exist node tops
						node := dot.NewNode(name)
						node.Set("group", dir)
						if _, ok := graphMap[parentNode.Name()]; !ok {
							graphMap[parentNode.Name()] = []string{}
						}
						edge := dot.NewEdge(parentNode, node)
						if label != "" {
							edge.Set("label", label)
						}
						switch dir {
						case dirPage:
							edge.Set("color", pageColor)
						case dirCon:
							edge.Set("color", contrColor)
						case dirBlock:
							edge.Set("color", pageColor)
						case dirMenu:
							edge.Set("color", menuColor)
						}
						graphDot.AddEdge(edge)
						graphMap[parentNode.Name()] = append(graphMap[parentNode.Name()], name)
					}
				}
			}
		}
	}
}

func getNodeName(name, dir string) (_name string) {
	_name = fmt.Sprintf("%s\n%s", name, strings.TrimSuffix(dir, "s"))
	if strings.Contains(_name, ",") {
		_name = strings.Join(strings.Split(_name, ","), "\n")
	}
	return
}
