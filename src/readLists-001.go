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
	"bufio"
	"encoding/csv"
	"io"
	"os"
)

// parse AntiBodies tag table
func parseAB(path string, ch2 chan<- map[string]string) {
	dataDict := make(map[string]string) // mutated AB tag => AB name (TCTCAGACCTCCGTA => CD14)

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		mutAB := mutSeq(record[0]) // compute all possible sequences with one mutation
		for _, ms := range mutAB {
			dataDict[ms] = record[1]
		}

	}
	ch2 <- dataDict
}

// parse white list cells
func parseWL(path string, ch1 chan<- CellRec) {
	dataDict := make(map[string]string) // mutated cell tag => cell "main" tag (AAACCTGCAATGGACG =>AAACCTGCAATGGACT)
	var cellsRefTag []string

	csvFile, err := os.Open(path)
	check(err)
	defer csvFile.Close()
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	for {
		// Read in a row. Check if we are at the end of the file.
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		cellsRefTag = append(cellsRefTag, record[0])
		mutS := mutSeq(record[0]) // compute all possible sequences with one mutation
		for _, ms := range mutS {
			dataDict[ms] = record[0]
		}

	}
	ch1 <- CellRec{dataDict, cellsRefTag}
}
