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

// isWhite search if cell sequence from cbf to cbl contain a cell tag in white list
func isWhite(seq string, cellsSeq map[string]string, conf CONF) (bool, string, string) {
	if len(seq) >= conf.Cbl {
		cellTag := seq[conf.Cbf-1 : conf.Cbl]
		if _, ok := cellsSeq[cellTag]; ok {
			return true, cellsSeq[cellTag], cellTag
		}
	}
	return false, "", ""
}

// isABody search if cell sequence from Abf to Abl contain a AB tag in tag list
func isABody(seq string, ABseq map[string]string, conf CONF) (bool, string) {
	if len(seq) >= conf.Abl {
		abTag := seq[conf.Abf-1 : conf.Abl]
		if abName, ok := ABseq[abTag]; ok {
			return true, abName
		}
	}
	return false, ""
}
