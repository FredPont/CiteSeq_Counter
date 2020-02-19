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
	"os"
	"sort"
	"strconv"
	"strings"
)

// return unique + sorted values of a map
func values(m map[string]string) []string {
	v := make([]string, 0, len(m))
	for _, value := range m {
		v = append(v, value)
	}
	v = uniqueStrings(v)
	sort.Strings(v)
	return v
}

// remove duplicates in []string
func uniqueStrings(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}

	return u
}

// write results table
func writeTable(abCounts map[CellAB]int, ABseq map[string]string, cellsRefTag []string) {

	// open result file for write
	fout := "result/result.tsv"
	out, err1 := os.Create(fout)
	check(err1)
	defer out.Close()

	ABslice := values(ABseq)
	header := "id\t" + strings.Join(ABslice, "\t") + "\n"
	writeOneLine(out, header) // write header in result file

	for _, c := range cellsRefTag {
		curRow := c + "\t"
		for _, ab := range ABslice {
			curRow = curRow + strconv.Itoa(abCounts[CellAB{c, ab}]) + "\t"
		}
		writeOneLine(out, curRow+"\n")
	}

}

//###########################################
func writeOneLine(f *os.File, line string) {
	_, err := f.WriteString(line)
	check(err)
}

//###########################################
