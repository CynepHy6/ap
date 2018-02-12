package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	currentVersion = "apla packager v0.7.3"

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
)

var (
	// flags
	condition      string
	menu           string
	outputName     string
	inputName      string
	permission     string
	unpackMode     bool
	verbose        bool
	version        bool
	singleSeparate bool
	dirs           []string
	// graphMode      bool
	sufMode bool
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

func init() {
	flag.BoolVar(&unpackMode, "unpack", false, "-u, unpacking mode")
	flag.StringVar(&inputName, "input", ".", "-i, path for input files, filename for pack and dirname/ (slashed) for unpack")
	flag.StringVar(&outputName, "output", "output", "-o, output filename for JSON if input file name not pointed")
	flag.StringVar(&condition, "conditions", "ContractConditions(\"MainCondition\")", "-c, conditions. Used if entry not founded in 'config.json'")
	flag.StringVar(&menu, "menu", "default_menu", "-m, menu. Used if entry not founded in 'config.json'")
	flag.StringVar(&permission, "table-permission", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "-t, permission for tables. Used if entry not founded in 'config.json'")
	flag.BoolVar(&verbose, "verbose", false, "print log")

	// shorthand
	flag.StringVar(&menu, "m", "default_menu", "-menu")
	flag.StringVar(&condition, "c", "ContractConditions(\"MainCondition\")", "--conditions")
	flag.StringVar(&outputName, "o", "output", "-output")
	flag.StringVar(&inputName, "i", ".", "input")
	flag.StringVar(&permission, "t", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "table-permission")
	flag.BoolVar(&unpackMode, "u", false, "-unpack")
	flag.BoolVar(&version, "v", false, "-version")
	flag.BoolVar(&singleSeparate, "s", false, "language and parameters will unpack to single separate files")
	// flag.BoolVar(&graphMode, "g", false, "-graph")
	// flag.BoolVar(&graphMode, "graph", false, "visualize call graph of package using dot format")
	flag.BoolVar(&sufMode, "suf", false, "unpack with suffixes for type")
	flag.Parse()

	dirs = []string{dirBlock, dirMenu, dirLang, dirTable, dirParam, dirData, dirPage, dirCon}
	if version {
		fmt.Println(currentVersion)
	}
	args := os.Args
	if len(args) == 1 {
		SimpleGui()
	} else {
		checkOutput()
	}
}

func main() {
	switch {
	case unpackMode:
		unpackJSON(inputName)
	default:
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
			fmt.Println("please choose file for unpaking, example:\n ap -u -i file.json")
			return //todo: create batch unpacking on Dir
		}
		if !strings.HasSuffix(outputName, separator) {
			outputName = outputName + separator
		}
		if verbose {
			fmt.Println("output dir name:", outputName)
		}
	} else {
		if !strings.HasSuffix(inputName, separator) {
			fmt.Println("please choose directory for paking, example:\n   ap -i dirfiles" + separator)
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
	if verbose {
		fmt.Println("extract:", outFile.Name())
	}

}
