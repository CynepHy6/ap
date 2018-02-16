package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func packJSON(path string) {
	initGraph()
	out := packDir(path)

	path = filepath.Dir(path)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}

	for _, f := range files {
		fname := f.Name()
		fpath := filepath.Join(path, fname)
		if debug {
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
			if debug {
				fmt.Println(err)
			}
			return
		}
		defer outFile.Close()
		outFile.WriteString(string(result))

		if abs, err := filepath.Abs(path); err == nil {
			abspath := filepath.Join(abs, structFileName)
			writeGraph(abspath)
		}
	}
	if debug {
		fmt.Println("not found files")
	}
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

	absdir, _ := filepath.Abs(path)
	absdirParts := strings.Split(absdir, separator)
	fdir := absdirParts[len(absdirParts)-1]
	for _, f := range files {
		fname := f.Name()
		ext := filepath.Ext(fname)
		name := fname[:len(fname)-len(ext)]
		if debug {
			fmt.Println(fname)
		}

		switch ext {
		case ePTL:
			switch {
			case strings.HasSuffix(name, _menu) || fdir == dirMenu:
				el := encodeStd(path, fname, _menu)
				createNodeForString(el.Name, fdir, el.Value)
				out.Menus = append(out.Menus, el)
			case strings.HasSuffix(name, _block) || fdir == dirBlock:
				el := encodeStd(path, fname, _block)
				createNodeForString(el.Name, fdir, el.Value)
				out.Blocks = append(out.Blocks, el)
			default:
				el := encodePage(path, fname, _page)
				createNodeForString(el.Name, fdir, el.Value)
				out.Pages = append(out.Pages, el)
			}
		case eJSON:
			switch {
			case name == "parameters":
				p := filepath.Join(path, fname)
				out.Parameters = append(out.Parameters, file2stdArray(p)...)
			case name == "languages":
				p := filepath.Join(path, fname)
				out.Languages = append(out.Languages, file2lang(p)...)
			case strings.HasSuffix(name, _param) || fdir == dirParam:
				out.Parameters = append(out.Parameters, encodeStd(path, fname, _param))
			case strings.HasSuffix(name, _lang) || fdir == dirLang:
				out.Languages = append(out.Languages, encodeLang(path, fname, _lang))
			case strings.HasSuffix(name, _table) || fdir == dirTable:
				out.Tables = append(out.Tables, encodeTable(path, fname, _table))
			case strings.HasSuffix(name, _data) || fdir == dirData:
				out.Data = append(out.Data, encodeData(path, fname, _data))
			}
		case eCSV:
			switch {
			case strings.HasSuffix(name, _param) || fdir == dirParam:
				out.Parameters = append(out.Parameters, encodeStd(path, fname, _param))
			}
		case eSIM:
			el := encodeStd(path, fname, _contr)
			createNodeForString(el.Name, fdir, el.Value)
			out.Contracts = append(out.Contracts, el)
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
func _JSONMarshal(v interface{}, unescape bool) ([]byte, error) {
	b, err := json.MarshalIndent(v, "", "    ")

	if unescape {
		b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
		b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
		b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
	}
	return b, err
}
func sortContracts(c []stdStruct) []stdStruct {
	nn := int(len(c) / 2)
	for n := 0; n < nn; n++ {
		for i := len(c) - 1; i > 0; i-- {
			for j := i - 1; j >= 0; j-- {
				if textContainsName(c[j].Value, c[i].Name) {
					c[i], c[j] = c[j], c[i]
					break
				}
			}
		}
	}
	return c
}

func textContainsName(text, name string) bool {
	lines := strings.Split(text, "\n")
	for _, l := range lines {
		line := strings.Trim(l, " ")
		if !strings.HasPrefix(line, "//") && strings.Contains(line, name) {
			return true
		}
	}
	return false
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
