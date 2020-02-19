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

// add a singe mutation at any position of a sequence
// ex :  mutSeq("TTT") -> [ATT CTT GTT TTT TAT TCT TGT TTT TTA TTC TTG TTT]

package main

func mutSeq(seq string) []string {
	nucl := []string{"A", "C", "G", "T"}
	var mutSeq []string
	for i := 0; i < len(seq); i++ {
		for _, n := range nucl {
			s := seq[:i] + n + seq[i+1:]
			mutSeq = append(mutSeq, s)
		}
	}
	return mutSeq
}
