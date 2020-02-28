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
	"fasta"
	"fmt"
	"math"
	"os"
	"regexp"

	"github.com/schollz/progressbar"
)

// BCumi = cell bar code + UMI used as key dictionnary to store ID
type BCumi struct {
	bc    string // cell bar code
	umi   string // umi
	mutBC string // cell bar code with possible mutation
}

type seenBCumi struct {
	BC string // cell bar code with no mutation
	AB string // AB name
}

// readR1R2 parses R1 and R2 file and filters reads with :
// - a valid cell tag between Cbf-1 : Cbl
// - a cell tag in the white list + 1 optionnal mutation
// - a valid AB tag from Abf to Abl ie a AB tag in tag list
// one mutation allowed in cell tag  + one mutation allowed inUMI
// reference cell tags (not mutated) + ABname + UMI used as reference
func readR1R2(fasta1, fasta2 string, conf CONF, cellsSeq map[string]string, ABseq map[string]string) map[CellAB]int { //
	abCounts := make(map[CellAB]int)       // (Cell+AB)=>count dictionnary to store AB counts
	seenRead := make(map[seenBCumi]string) // bar code + AB already counted => umi
	unvalidRead := 0                       // read with valid AB but not valid cell
	PCRduplicates := 0                     // PCR duplicates

	reg, _ := regexp.Compile(conf.Treg) // compile TAG regex "^[ATGC]{15}[TGC][A]{6,}"

	// measure file size
	nbSeq := SeqEstim(fasta1)
	fmt.Println("filename =  ", fasta1, " size = ", nbSeq, "sequences\n")

	// Open the files.
	f1, err := os.Open(fasta1) // open fasta file R1
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	f2, err := os.Open(fasta2) // open fasta file R2
	if err != nil {
		panic(err)
	}
	defer f2.Close()

	// Loop over all lines in the fasta file
	var fqr1 fasta.FqReader
	re1 := bufio.NewReader(f1)
	fqr1.Reader = re1

	var fqr2 fasta.FqReader
	re2 := bufio.NewReader(f2)
	fqr2.Reader = re2
	n := 0

	//var test []string
	count := int(nbSeq)
	bar := progressbar.New(count) // Add a new progress bar
	for {
		bar.Add(1) // show progress bar

		re1, done1 := fqr1.Iter()
		re2, done2 := fqr2.Iter()

		if done1 || done2 {
			break
		}

		n++

		inlist, abName := isABody(re2.Seq, ABseq, conf)
		// if cell AB TAG is valid
		if !inlist {
			continue
		} else {
			//inlist, cellTag, cellMutBC := isWhite(re1.Seq, cellsSeq, conf)
			inlist, cellTag, _ := isWhite(re1.Seq, cellsSeq, conf)
			// if cell ID TAG, and regex are good and if cell is in not in white list
			if !inlist {
				unvalidRead++
				continue
			} else {
				if reg.MatchString(re2.Seq) { // regex is slow it is faster to check it last
					umi := re1.Seq[conf.Umif-1 : conf.Umil] // UMI
					//fmt.Println(cellTag, cellMutBC, abName, umi)

					// one mutation allowed in cell tag  + one mutation allowed inUMI
					// reference cell tags (not mutated) + ABname + UMI used as reference

					seenUMI, ok := seenRead[seenBCumi{cellTag, abName}] // test if bc + AB + umi was already counted ie PCR duplicate
					if ok {
						if umi == seenUMI || ld(umi, seenUMI) < 2 {
							PCRduplicates++
							continue // if the UMI was seen, it is a PCR duplicate nothing is counted
						} else {
							ct := abCounts[CellAB{cellTag, abName}]
							abCounts[CellAB{cellTag, abName}] = ct + 1
						}

					} else { // if this  barcode + AB name is unknown, it is recorded
						seenRead[seenBCumi{cellTag, abName}] = umi

						abCounts[CellAB{cellTag, abName}] = 1

					}

				}
			}
		}

	}
	fmt.Println(unvalidRead, " AB reads rejected\n")
	fmt.Println(PCRduplicates, " PCR duplicates\n")
	fmt.Println(n, " sequences read\n")

	//fmt.Println(len(test), " AB valid sequences read\n")
	return abCounts
}

// SeqEstim estimate number of fasta sequence from the file size
func SeqEstim(filename string) float64 {
	var nbSeq float64
	fi, e := os.Stat(filename)
	if e != nil {
		fmt.Println("can't read ", filename, "size, ", e)
	}
	// get the size
	size := float64(fi.Size())
	// calculate the nb of sequence
	// R1 : 510682384 for 14112648 Kbytes
	// R2 : 510682384 for 32066328 Kbytes

	nbSeq = math.Ceil(510682384.0 * size * 0.001 * .25 / 14112648.0) // *.25 : 4 lignes/seq ,  * 0.001 kbytes -> bytes

	return nbSeq
}
