package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type configFile struct {
	Blocks    *[]stdConf   `json:"blocks"`
	Contracts *[]stdConf   `json:"contracts"`
	Menus     *[]stdConf   `json:"menus"`
	Pages     *[]pageConf  `json:"pages"`
	Tables    *[]tableConf `json:"tables"`
	Params    *[]stdConf   `json:"parameters"`
}
type stdConf struct {
	Name       string
	Conditions string
}

type pageConf struct {
	Name       string
	Conditions string
	Menu       string
}

type tableConf struct {
	Name        string
	Permissions string
}

type exportFile struct {
	Blocks     []stdStruct   `json:"blocks"`
	Contracts  []stdStruct   `json:"contracts"`
	Data       []dataStruct  `json:"data"`
	Languages  []langStruct  `json:"languages"`
	Menus      []stdStruct   `json:"menus"`
	Pages      []pageStruct  `json:"pages"`
	Parameters []stdStruct   `json:"parameters"`
	Tables     []tableStruct `json:"tables"`
}
type importFile struct {
	Blocks     []commonStruct `json:"blocks"`
	Contracts  []commonStruct `json:"contracts"`
	Data       []dataStruct   `json:"data"`
	Languages  []commonStruct `json:"languages"`
	Menus      []commonStruct `json:"menus"`
	Pages      []commonStruct `json:"pages"`
	Parameters []commonStruct `json:"parameters"`
	Tables     []commonStruct `json:"tables"`
}

type commonStruct struct {
	Name       string
	Value      string
	Conditions string
	Trans      string
	Columns    string
	Table      string
}
type stdStruct struct {
	Name       string
	Value      string
	Conditions string
}
type langStruct struct {
	Name       string
	Conditions string
	Trans      string
}

type pageStruct struct {
	Name       string
	Value      string
	Conditions string
	Menu       string
}

type tableStruct struct {
	Name        string
	Columns     string
	Permissions string
}

type dataStruct struct {
	Table   string
	Columns []string
	Data    [][]string
}

const (
	currentVersion = "apla packager v0.8.2"

	eSIM  = ".sim"
	ePTL  = ".ptl"
	eJSON = ".json"
	eCSV  = ".csv"

	// file suffixes
	_block = "__block"
	_menu  = "__menu"
	_lang  = "__language"
	_table = "__table"
	_param = "__parameter"
	_data  = "__data"
	_page  = "__page"
	_contr = "__contract"

	//dirs
	dirBlock = "blocks"
	dirMenu  = "menus"
	dirLang  = "languages"
	dirTable = "tables"
	dirParam = "parameters"
	dirData  = "data"
	dirPage  = "pages"
	dirCon   = "contracts"

	//
	configName = "config.json"
	separator  = string(os.PathSeparator)

	structFileName = "struct.dot"
	pageColor      = "green"
	contrColor     = "red"
	menuColor      = "blue"

	// help messages
	helpMsg = "please choose directory for paking, example:\n    ap dirfiles" + separator + "\nor file to unpacking, example:\n    ap file.json"
)

var (
	// flags
	condition      = "ContractConditions(\"MainCondition\")"
	menu           = "default_menu"
	outputName     string
	inputName      string
	permission     = "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}"
	unpackMode     bool
	debug          bool
	version        bool
	singleSeparate bool
	// graphMode      bool
	sufMode bool
	dirs    = []string{dirBlock, dirMenu, dirLang, dirTable, dirParam, dirData, dirPage, dirCon}
)

func init() {
	flag.BoolVar(&unpackMode, "unpack", false, "-u, unpacking mode")
	flag.StringVar(&inputName, "input", ".", "-i, path for input files, filename for pack and dirname/ (slashed) for unpack")
	flag.StringVar(&outputName, "output", "output", "-o, output filename for JSON if input file name not pointed")

	// shorthand
	flag.StringVar(&outputName, "o", "output", "-output")
	flag.StringVar(&inputName, "i", ".", "input")
	flag.BoolVar(&unpackMode, "u", false, "-unpack")
	flag.BoolVar(&version, "v", false, "-version")
	flag.BoolVar(&debug, "d", false, "debug")
	flag.Parse()
}

func main() {
	args := os.Args

	if argsCount := len(args); argsCount == 1 {
		// without args run gui
		SimpleGui()
	} else {
		if argsCount == 2 {
			if version {
				fmt.Println(currentVersion)
			} else {
				name := args[1]
				inputName = name
				if !strings.HasSuffix(name, separator) { // if filename run unpack
					unpackMode = true
				}
			}
		} else {
			if version {
				fmt.Println(currentVersion)
			}
		}
		checkOutput()
	}
	if unpackMode {
		unpackJSON(inputName)
	} else {
		packJSON(inputName)
	}
}

func checkOutput() {
	if outputName == "output" && inputName != "." { // we have only inputname
		parts := strings.Split(inputName, separator)
		pLen := len(parts)
		outputName = parts[pLen-1]
		if unpackMode {
			ext := filepath.Ext(outputName)
			outputName = outputName[:len(outputName)-len(ext)]
			outputName = outputName + separator
		} else {
			if strings.HasSuffix(inputName, separator) {
				outputName = parts[pLen-2]
			}
		}
	}

	if unpackMode {
		if stats, err := os.Stat(inputName); inputName == "." || stats.IsDir() || err != nil {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(helpMsg)
			return
		}
		if !strings.HasSuffix(outputName, separator) {
			outputName = outputName + separator
		}
		if debug {
			fmt.Println("output dir name:", outputName)
		}
	} else {
		if !strings.HasSuffix(inputName, separator) {
			fmt.Println(helpMsg)
			return
		}
	}
}

func createDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		fmt.Println("error create dir", err)
	}
}

func writeFileString(filename, content string) {
	outFile, err := os.Create(filepath.Join(outputName, filename))
	if err != nil {
		// fmt.Println("error write file:", err)
		return
	}
	defer outFile.Close()
	if _, err := outFile.WriteString(content); err != nil {
		fmt.Println(err)
		return
	}
}

func stringInSlice(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
