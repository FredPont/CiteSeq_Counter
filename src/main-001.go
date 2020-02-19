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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// CellAB = cell name + antibody name used as key dictionnary to store AB counts
type CellAB struct {
	cell string // cell tag
	ab   string // antibody tag
}

// CellRec is used in the channel ch2 to retrieve cells TAG map + cells ref tags
type CellRec struct {
	ctm map[string]string // cells TAG map
	rt  []string          // ref tags
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// ListFiles lists all files in a directory
func ListFiles(dir string) []string {
	var filesList []string
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println(f.Name(), "in", dir)
		filesList = append(filesList, f.Name())
	}
	return filesList
}

// checkFileNB verify that there is only one file in the file list
func checkFileNB(files []string) {
	if len(files) > 1 {
		nbTables := strconv.Itoa(len(files))
		msg := "error ! " + nbTables + " tag files found, only one permitted ! : " + strings.Join(files, " ")
		err := errors.New(msg)
		check(err)
	}
}

func header() {

	fmt.Println("   ┌──────────────────────────────────────────┐") // unicode U+250C
	fmt.Println("   │  CITE-seq-counter (c)Frederic PONT 2018  │")
	fmt.Println("   │ Free Software GNU General Public License │")
	fmt.Println("   └──────────────────────────────────────────┘")
}

// important variables summary
// cellsSeq : mutated cell tag => cell "main" tag (AAACCTGCAATGGACG => AAACCTGCAATGGACT)
// cellsMap : id -> "main" barcode (@NB551452:13:H52L5BGX9:1:21201:3572:11047 => CTCATTGTAACTCCT)
// ABseq 	: mutated AB tag => AB name (TCTCAGACCTCCGTA => CD14)

func main() {
	header()
	t0 := time.Now()

	// read configuration file
	conf := ReadConfig()
	fmt.Println("Config summary :", conf)

	// read white list cells
	whiteL := ListFiles("whiteList")
	checkFileNB(whiteL)
	path := "whiteList/" + whiteL[0]
	ch1 := make(chan CellRec)
	go parseWL(path, ch1)

	// read antibody signatures
	ABtags := ListFiles("tags")
	checkFileNB(ABtags)
	path = "tags/" + ABtags[0]
	ch2 := make(chan map[string]string)
	go parseAB(path, ch2)

	ABseq := <-ch2 // mutated AB tag => AB name (TCTCAGACCTCCGTA => CD14)
	fmt.Println(len(ABseq), "possible AB sequences in tags list")

	record := <-ch1 // mutated cell tag => cell "main" tag (AAACCTGCAATGGACG =>AAACCTGCAATGGACT)
	cellsSeq := record.ctm
	cellsRefTag := record.rt
	fmt.Println(len(cellsSeq), "possible cell sequences in white list")

	// read and filter R1
	cellsReads := ListFiles("fastqR1")
	checkFileNB(cellsReads)
	path1 := "fastqR1/" + cellsReads[0]

	// read and filter R2
	cellsReads = ListFiles("fastqR2")
	checkFileNB(cellsReads)
	path2 := "fastqR2/" + cellsReads[0]
	fmt.Println("R1 and R2 parsing...")
	abCounts := readR1R2(path1, path2, conf, cellsSeq, ABseq)

	writeTable(abCounts, ABseq, cellsRefTag)

	fmt.Println("Finished !")
	fmt.Printf("Elapsed time : %v.\n", time.Since(t0))
	fmt.Print("Press enter to close window ")
	//fmt.Scanln() // saisie clavier
}
