package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	currentVersion = "apla packager v0.6.6"

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
	unpack         bool
	verbose        bool
	version        bool
	singleSeparate bool
	dirs           []string
	graph          bool
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
	flag.BoolVar(&unpack, "-unpack", false, "-u, unpacking mode")
	flag.StringVar(&inputName, "-input", ".", "-i, path for input files, filename for pack and dirname/ (slashed) for unpack")
	flag.StringVar(&outputName, "-output", "output", "-o, output filename for JSON if input file name not pointed")
	flag.StringVar(&condition, "-conditions", "ContractConditions(\"MainCondition\")", "-c, conditions. Used if entry not founded in 'config.json'")
	flag.StringVar(&menu, "-menu", "default_menu", "-m, menu. Used if entry not founded in 'config.json'")
	flag.StringVar(&permission, "-table-permission", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "-t, permission for tables. Used if entry not founded in 'config.json'")
	flag.BoolVar(&verbose, "-verbose", false, "print log")

	// shorthand
	flag.StringVar(&menu, "m", "default_menu", "--menu")
	flag.StringVar(&condition, "c", "ContractConditions(\"MainCondition\")", "--conditions")
	flag.StringVar(&outputName, "o", "output", "--output")
	flag.StringVar(&inputName, "i", ".", "--input")
	flag.StringVar(&permission, "t", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "--table-permission")
	flag.BoolVar(&unpack, "u", false, "--unpack")
	flag.BoolVar(&version, "v", false, "version")
	flag.BoolVar(&singleSeparate, "s", false, "language and parameters will unpack to single separate files")
	flag.BoolVar(&graph, "g", false, "visualize call graph of package using dot format")
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
	if unpack {
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
		if unpack {
			ext := filepath.Ext(outputName)
			outputName = outputName[:len(outputName)-len(ext)]
			outputName = outputName + separator
		} else {
			if strings.HasSuffix(inputName, separator) {
				outputName = parts[pLen-2]
			}
		}
	}

	if unpack {
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

func packJSON(path string) {

	out := packDir(path)

	path = filepath.Dir(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		fname := f.Name()
		fpath := filepath.Join(path, fname)
		if verbose {
			fmt.Println(fpath)
		}
		sf, err := os.Stat(fpath)
		if err != nil {
			fmt.Println(err)
			return
		}
		if sf.IsDir() {
			dir := packDir(fpath)
			switch fname {
			case dirBlock:
				out.Blocks = append(out.Blocks, dir.Blocks...)
			case dirMenu:
				out.Menus = append(out.Menus, dir.Menus...)
			case dirLang:
				out.Languages = append(out.Languages, dir.Languages...)
			case dirTable:
				out.Tables = append(out.Tables, dir.Tables...)
			case dirParam:
				out.Parameters = append(out.Parameters, dir.Parameters...)
			case dirData:
				out.Data = append(out.Data, dir.Data...)
			case dirPage:
				out.Pages = append(out.Pages, dir.Pages...)
			case dirCon:
				out.Contracts = append(out.Contracts, dir.Contracts...)
			}
		}
	}
	if countEntries(out) > 0 {
		readConfig(&out)
		if len(out.Contracts) > 0 {
			out.Contracts = sortContracts(out.Contracts)
		}

		result, _ := _JSONMarshal(out, true)
		if !strings.HasSuffix(outputName, ".json") {
			outputName += ".json"
		}
		outFile, err := os.Create(outputName)
		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			return
		}
		defer outFile.Close()
		outFile.WriteString(string(result))
	}
	if verbose {
		fmt.Println("not found files")
	}
}

func _JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}

func packDir(path string) (out exportFile) {
	out.Blocks = []stdStruct{}
	out.Contracts = []stdStruct{}
	out.Data = []dataStruct{}
	out.Languages = []langStruct{}
	out.Menus = []stdStruct{}
	out.Pages = []pageStruct{}
	out.Parameters = []stdStruct{}
	out.Tables = []tableStruct{}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		fname := f.Name()
		ext := filepath.Ext(fname)
		name := fname[:len(fname)-len(ext)]
		if verbose {
			fmt.Println(fname)
		}

		switch ext {
		case ePTL:
			switch {
			case strings.HasSuffix(name, _menu):
				out.Menus = append(out.Menus, encodeStd(path, fname, _menu))
			case strings.HasSuffix(name, _block):
				out.Blocks = append(out.Blocks, encodeStd(path, fname, _block))
			default:
				out.Pages = append(out.Pages, encodePage(path, fname, _page))
			}
		case eJSON:
			switch {
			case name == "parameters":
				p := filepath.Join(path, fname)
				out.Parameters = append(out.Parameters, file2stdArray(p)...)
			case name == "languages":
				p := filepath.Join(path, fname)
				out.Languages = append(out.Languages, file2lang(p)...)
			case strings.HasSuffix(name, _param):
				out.Parameters = append(out.Parameters, encodeStd(path, fname, _param))
			case strings.HasSuffix(name, _lang):
				out.Languages = append(out.Languages, encodeLang(path, fname, _lang))
			case strings.HasSuffix(name, _table):
				out.Tables = append(out.Tables, encodeTable(path, fname, _table))
			case strings.HasSuffix(name, _data):
				out.Data = append(out.Data, encodeData(path, fname, _data))
			}
		case eCSV:
			switch {
			case strings.HasSuffix(name, _param):
				out.Parameters = append(out.Parameters, encodeStd(path, fname, _param))
			}
		case eSIM:
			out.Contracts = append(out.Contracts, encodeStd(path, fname, _contr))
		}

	}
	return
}

func encodePage(path, fname, sExt string) (result pageStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}
	result.Menu = menu
	result.Name = name
	result.Value = file2str(fpath)
	result.Conditions = condition
	return
}
func encodeData(path, fname, sExt string) (result dataStruct) {
	// result = make(map[string]string)
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}
	result.Table = name
	dataFile := file2data(fpath)
	result.Columns = dataFile.Columns
	result.Data = dataFile.Data
	return
}
func encodeTable(path, fname, sExt string) (result tableStruct) {
	// result = make(map[string]string)
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}
	result.Name = name
	result.Columns = file2str(fpath)
	result.Permissions = permission
	return
}
func encodeLang(path, fname, sExt string) (result langStruct) {
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}
	result.Name = name
	result.Trans = file2str(fpath)
	result.Conditions = ""
	return
}

func encodeStd(path, fname, sExt string) (result stdStruct) {
	// result = make(map[string]string)
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}
	result.Name = name
	result.Value = file2str(fpath)
	result.Conditions = condition
	return
}

func file2str(filename string) (str string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	str = string(bs)
	return
}

func file2data(filename string) (result dataStruct) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(bs, &result)
	return
}

func file2stdArray(filename string) (result []stdStruct) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(bs, &result)
	return
}

func file2lang(filename string) (result []langStruct) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(bs, &result)
	return
}

func unpackJSON(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := exportFile{}
	if err := json.Unmarshal(bs, &file); err != nil {
		fmt.Println("unmarshal file error:", err)
		return
	}

	if len(file.Contracts) > 0 {
		createDir(filepath.Join(outputName, dirCon))
		for _, c := range file.Contracts {
			value := c.Value
			name := c.Name + eSIM
			name = filepath.Join(dirCon, name)
			writeFileString(name, value)
		}
	}
	if len(file.Menus) > 0 {
		createDir(filepath.Join(outputName, dirMenu))
		for _, c := range file.Menus {
			value := c.Value
			name := c.Name + _menu + ePTL
			name = filepath.Join(dirMenu, name)
			writeFileString(name, value)
		}
	}
	if len(file.Parameters) > 0 {
		if singleSeparate {
			byteValue, _ := json.MarshalIndent(file.Parameters, "", "    ")
			value := string(byteValue)
			name := dirParam + eJSON
			writeFileString(name, value)
		} else {
			createDir(filepath.Join(outputName, dirParam))
			for _, c := range file.Parameters {
				value := c.Value
				name := c.Name + _param + eCSV
				name = filepath.Join(dirParam, name)
				writeFileString(name, value)
			}
		}
	}
	if len(file.Languages) > 0 {
		if singleSeparate {
			byteValue, _ := json.MarshalIndent(file.Languages, "", "    ")
			value := string(byteValue)
			name := dirLang + eJSON
			writeFileString(name, value)
		} else {
			createDir(filepath.Join(outputName, dirLang))
			for _, c := range file.Languages {
				value := c.Trans
				name := c.Name + _lang + eJSON
				name = filepath.Join(dirLang, name)
				writeFileString(name, value)
			}
		}
	}
	if len(file.Tables) > 0 {
		createDir(filepath.Join(outputName, dirTable))
		for _, c := range file.Tables {
			value := c.Columns
			name := c.Name + _table + eJSON
			name = filepath.Join(dirTable, name)
			writeFileString(name, value)
		}
	}
	if len(file.Blocks) > 0 {
		createDir(filepath.Join(outputName, dirBlock))
		for _, c := range file.Blocks {
			value := c.Value
			name := c.Name + _block + ePTL
			name = filepath.Join(dirBlock, name)
			writeFileString(name, value)
		}
	}
	if len(file.Data) > 0 {
		createDir(filepath.Join(outputName, dirData))
		for _, c := range file.Data {
			name := c.Table + _data + eJSON
			name = filepath.Join(dirData, name)
			result, _ := json.MarshalIndent(c, "", "    ")
			writeFileString(name, string(result))
		}
	}
	if len(file.Pages) > 0 {
		createDir(filepath.Join(outputName, dirPage))
		for _, c := range file.Pages {
			value := c.Value
			name := c.Name + ePTL
			name = filepath.Join(dirPage, name)
			writeFileString(name, value)
		}
	}
	writeConfig(bs)
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
func readConfig(out *exportFile) {
	config := exportFile{}
	absConfPath, _ := filepath.Abs(inputName)
	bs, err := ioutil.ReadFile(filepath.Join(absConfPath, configName))
	if err != nil {
		if verbose {
			fmt.Println("config file not found. use default values")
		}
		return
	}
	_ = json.Unmarshal(bs, &config)
	if len(config.Blocks) > 0 {
		for c := range config.Blocks {
			for o := range out.Blocks {
				if config.Blocks[c].Name == out.Blocks[o].Name {
					out.Blocks[o].Conditions = config.Blocks[c].Conditions
				}
			}
		}
	}
	if len(config.Contracts) > 0 {
		for c := range config.Contracts {
			for o := range out.Contracts {
				if config.Contracts[c].Name == out.Contracts[o].Name {
					out.Contracts[o].Conditions = config.Contracts[c].Conditions
				}
			}
		}
	}
	if len(config.Menus) > 0 {
		for c := range config.Menus {
			for o := range out.Menus {
				if config.Menus[c].Name == out.Menus[o].Name {
					out.Menus[o].Conditions = config.Menus[c].Conditions
				}
			}
		}
	}
	if len(config.Pages) > 0 {
		for c := range config.Pages {
			for o := range out.Pages {
				if config.Pages[c].Name == out.Pages[o].Name {
					out.Pages[o].Conditions = config.Pages[c].Conditions
					if len(config.Pages[c].Menu) > 0 {
						out.Pages[o].Menu = config.Pages[c].Menu
					}
				}
			}
		}
	}
	if len(config.Tables) > 0 {
		for c := range config.Tables {
			for o := range out.Tables {
				if config.Tables[c].Name == out.Tables[o].Name {
					out.Tables[o].Permissions = config.Tables[c].Permissions
				}
			}
		}
	}
	if len(config.Parameters) > 0 {
		for c := range config.Parameters {
			for o := range out.Parameters {
				if config.Parameters[c].Name == out.Parameters[o].Name {
					out.Parameters[o].Conditions = config.Parameters[c].Conditions
				}
			}
		}
	}
	return
}
func writeConfig(bs []byte) {
	cFile := configFile{}
	if err := json.Unmarshal(bs, &cFile); err != nil {
		fmt.Println("unmarshal config file error:", err)
	} else {
		if bs, err := json.MarshalIndent(cFile, "", "    "); err == nil {
			writeFileString(configName, string(bs))
		}
	}
}
func sortContracts(c []stdStruct) (res []stdStruct) {
	// sorting contracts by used in other contracts
	res = c
	lenC := len(c)
	for i := lenC - 1; i > 0; i-- {
		name := c[i].Name
		for j := i - 1; j >= 0; j-- {
			value := c[j].Value
			if strings.Contains(value, name) {
				c[i], c[j] = c[j], c[i]
			}
		}
	}
	return
}

func countEntries(file exportFile) (count int) {
	return len(file.Blocks) +
		len(file.Contracts) +
		len(file.Data) +
		len(file.Languages) +
		len(file.Menus) +
		len(file.Pages) +
		len(file.Parameters) +
		len(file.Tables)
}
