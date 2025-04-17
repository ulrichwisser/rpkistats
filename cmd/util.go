/*
Copyright © 2025 Ulrich Wisser

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (

)


func unique(strlist []string) (resultlist []string) {
	resultlist = make([]string, 0)
	strmap := make(map[string]bool, 0)

	for _,s := range strlist {
		strmap[s] = true
	}

	for s := range strmap {
		resultlist = append(resultlist, s)
	}

	return
}