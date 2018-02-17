package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func readConfig(out *exportFile) {
	config := exportFile{}
	absConfPath, _ := filepath.Abs(inputName)
	bs, err := ioutil.ReadFile(filepath.Join(absConfPath, configName))
	if err != nil {
		if debug {
			fmt.Println("config file not found. used default values")
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
