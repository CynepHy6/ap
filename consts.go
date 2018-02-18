package main

import "os"

const (
	currentVersion = "apla packager 0.8.3"

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

	dirBlock = "blocks"
	dirMenu  = "menus"
	dirLang  = "languages"
	dirTable = "tables"
	dirParam = "parameters"
	dirData  = "data"
	dirPage  = "pages"
	dirCon   = "contracts"

	defaultCondition  = "ContractConditions(\"MainCondition\")"
	defaultMenu       = "default_menu"
	defaultPermission = "{\"insert\":\"true\",\"update\":\"true\",\"new_column\":\"true\"}"
	configName        = "config.json"
	separator         = string(os.PathSeparator)
	structFileName    = "struct.dot"

	//
	helpMsg = "please choose directory for paking, example:\n    ap dirfiles" + separator + "\nor file to unpacking, example:\n    ap file.json"
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

type graphStruct struct {
	Name      string
	Value     string
	Group     string
	Path      string
	Dir       string
	FontColor string
	Color     string
	EdgeLabel string
}
