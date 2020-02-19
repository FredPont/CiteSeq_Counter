/*
 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program.  If not, see <http://www.gnu.org/licenses/>.

 Written by Frederic PONT.
 (c) Frederic Pont 2018
*/

package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// CONF stores software parameters
type CONF struct {
	Cbf  int    `json:"cell_barcode_first_base"`
	Cbl  int    `json:"cell_barcode_last_base"`
	Umif int    `json:"umi_first_base"`
	Umil int    `json:"umi_last_base"`
	Abf  int    `json:"AB_barcode_first_base"`
	Abl  int    `json:"AB_barcode_last_base"`
	Treg string `json:"TAG_regex"`
}

// ReadConfig reads conf.json
func ReadConfig() CONF {
	file, err1 := os.Open("conf.json")
	if err1 != nil {
		fmt.Println(err1)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := CONF{}
	err := decoder.Decode(&conf)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(conf)
	return conf
}
