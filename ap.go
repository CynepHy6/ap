package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	eSIM  = ".sim"
	ePTL  = ".ptl"
	eJSON = ".json"

	// file suffixes
	_block = "__block"
	_menu  = "__menu"
	_lang  = "__language"
	_table = "__table"
	_param = "__parameter"
	_data  = "__data"
	_page  = "__page"
	_contr = "__contract"
)

var (
	// flags
	condition  string
	menu       string
	outName    string
	prefix     string
	path       string
	permission string
	unpack     bool
	verbose    bool
)

func init() {
	flag.StringVar(&condition, "c", "ContractConditions(\"MainCondition\")", "--conditions")
	flag.StringVar(&condition, "-conditions", "ContractConditions(\"MainCondition\")", "-c, conditions")
	flag.StringVar(&menu, "m", "default_menu", "--menu")
	flag.StringVar(&menu, "-menu", "default_menu", "-m, menu")
	flag.StringVar(&outName, "o", "output", "--output")
	flag.StringVar(&outName, "-output", "output", "-o, output filename for JSON")
	flag.StringVar(&prefix, "p", "", "--prefix")
	flag.StringVar(&prefix, "-prefix", "", "-p, prefix for pages and contracts")
	flag.StringVar(&path, "i", ".", "--input")
	flag.StringVar(&path, "-input", ".", "-i, path for input files")
	flag.StringVar(&permission, "t", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "--table-permission")
	flag.StringVar(&permission, "-table-permission", "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}", "-t, permission for tables")
	flag.BoolVar(&unpack, "-unpack", false, "-u, unpacking mode")
	flag.BoolVar(&unpack, "u", false, "--unpack")
	flag.BoolVar(&verbose, "-verbose", false, "work log")
	flag.BoolVar(&verbose, "v", false, "--verbose")
}

func main() {
	flag.Parse()
	if prefix != "" {
		prefix = prefix + "_"
		outName = prefix + outName
	}
	if unpack {
		if stats, err := os.Stat(path); path == "." || stats.IsDir() || err != nil {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("please choose file for unpaking, example:\n ap -u -i file.json")
			return //todo: create batch unpacking on Dir
		}
		if !strings.HasSuffix(outName, string(os.PathSeparator)) {
			outName = outName + string(os.PathSeparator)
		}
		if verbose {
			fmt.Println("output dir name:", outName)
		}
		unpackJSON(path)

	} else {
		content := packJSON(path)
		if content == "" {
			return
		}
		outFile, err := os.Create(outName + ".json")
		if err != nil {
			if verbose {
				fmt.Println(err)
			}
			return
		}
		defer outFile.Close()
		outFile.WriteString(content)
	}
}

func packJSON(path string) string {
	out := make(map[string][]map[string]string)
	emptyMap := []map[string]string{}
	contracts := emptyMap
	pages := emptyMap
	menus := emptyMap
	params := emptyMap
	langs := emptyMap
	tables := emptyMap
	datas := emptyMap
	blocks := emptyMap
	path = filepath.Dir(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	var countFiles int
	for _, f := range files {
		fname := f.Name()
		ext := filepath.Ext(fname)
		name := fname[:len(fname)-len(ext)]
		if verbose {
			fmt.Println(fname)
		}
		switch ext {
		case ePTL:
			countFiles++
			switch {
			case strings.HasSuffix(name, _menu):
				menus = append(menus, encode(path, fname, _menu))
			case strings.HasSuffix(name, _block):
				blocks = append(blocks, encode(path, fname, _block))
			default:
				pages = append(pages, encode(path, fname, _page))
			}
		case eJSON:
			switch {
			case strings.HasSuffix(name, _param):
				countFiles++
				params = append(params, encode(path, fname, _param))
			case strings.HasSuffix(name, _lang):
				countFiles++
				langs = append(langs, encode(path, fname, _lang))
			case strings.HasSuffix(name, _table):
				countFiles++
				tables = append(tables, encode(path, fname, _table))
			case strings.HasSuffix(name, _data):
				countFiles++
				datas = append(datas, encode(path, fname, _data))
			}
		case eSIM:
			countFiles++
			contracts = append(contracts, encode(path, fname, _contr))
		}
	}
	if countFiles > 0 {
		out["menus"] = menus
		out["parameters"] = params
		out["languages"] = langs
		out["tables"] = tables
		out["data"] = datas
		out["blocks"] = blocks
		out["pages"] = pages
		out["contracts"] = contracts
		result, _ := json.Marshal(out)
		return string(result)
	}
	if verbose {
		fmt.Println("not found files")
	}
	return ""
}

func encode(path, fname, sExt string) (result map[string]string) {
	result = make(map[string]string)
	ext := filepath.Ext(fname)
	name := fname[:len(fname)-len(ext)]
	fpath := filepath.Join(path, fname)
	if strings.HasSuffix(name, sExt) {
		// remove suffix
		name = name[:len(name)-len(sExt)]
	}

	switch sExt {
	case _contr:
		result["Name"] = prefix + name
		result["Value"] = file2str(fpath)
		if prefix != "" {
			// apply prefix in contract on name
			re := regexp.MustCompile("contract\\s+" + name)
			result["Value"] = re.ReplaceAllString(result["Value"], "contract "+result["Name"])
		}
		result["Conditions"] = condition
	case _menu:
		result["Name"] = prefix + name
		result["Value"] = file2str(fpath)
		result["Conditions"] = condition
	case _param:
		result["Name"] = prefix + name
		result["Value"] = file2str(fpath)
		result["Conditions"] = condition
	case _lang:
		result["Name"] = prefix + name
		result["Trans"] = file2str(fpath)
		result["Conditions"] = ""
	case _table:
		result["Name"] = prefix + name
		result["Columns"] = file2str(fpath)
		result["Permissions"] = permission
	case _block:
		result["Name"] = prefix + name
		result["Value"] = file2str(fpath)
		result["Conditions"] = condition
	case _data:
		result["Table"] = prefix + name
		dataTable := file2data(fpath)
		result["Columns"] = dataTable["Columns"]
		result["Data"] = dataTable["Data"]
	case _page:
		result["Menu"] = menu
		result["Name"] = prefix + name
		result["Value"] = file2str(fpath)
		result["Conditions"] = condition
	}
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

func file2data(filename string) (result map[string]string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	json.Unmarshal(bs, result)
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
		fmt.Println(err)
		return
	}
	if len(file.Contracts) > 0 {
		for _, c := range file.Contracts {
			value := c.Value
			name := prefix + c.Name + eSIM
			writeFileString(name, value)
		}
	}
	if len(file.Menus) > 0 {
		for _, c := range file.Menus {
			value := c.Value
			name := prefix + c.Name + _menu + ePTL
			writeFileString(name, value)
		}
	}
	if len(file.Parameters) > 0 {
		for _, c := range file.Parameters {
			value := c.Value
			name := prefix + c.Name + _param + eJSON
			writeFileString(name, value)
		}
	}
	if len(file.Languages) > 0 {
		for _, c := range file.Languages {
			value := c.Trans
			name := prefix + c.Name + _lang + eJSON
			writeFileString(name, value)
		}
	}
	if len(file.Tables) > 0 {
		for _, c := range file.Tables {
			value := c.Columns
			name := prefix + c.Name + _table + eJSON
			writeFileString(name, value)
		}
	}
	if len(file.Blocks) > 0 {
		for _, c := range file.Blocks {
			value := c.Value
			name := prefix + c.Name + _block + ePTL
			writeFileString(name, value)
		}
	}
	if len(file.Data) > 0 {
		for _, c := range file.Data {
			name := prefix + c.Table + _data + eJSON
			outFile, err := os.Create(filepath.Join(outName, name))
			if err != nil {
				continue
			}
			defer outFile.Close()
			result, _ := json.Marshal(c)
			writeFileString(name, string(result))
		}
	}
	if len(file.Pages) > 0 {
		for _, c := range file.Pages {
			value := c.Value
			name := prefix + c.Name + ePTL
			writeFileString(name, value)
		}
	}
}

func writeFileString(filename, content string) {
	if err := os.MkdirAll(outName, os.ModePerm); err != nil {
		fmt.Println(err)
	}
	outFile, err := os.Create(filepath.Join(outName, filename))
	if err != nil {
		fmt.Println(err)
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
