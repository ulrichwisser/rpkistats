/*
Copyright © 2023 Ulrich Wisser <ulrich@wisser.se>

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
	"github.com/spf13/viper"
)

const VERBOSE = "verbose"
const VERBOSE_SHORT = "v"
const VERBOSE_QUIET int = 0
const VERBOSE_ERROR int = 1
const VERBOSE_WARNING int = 2
const VERBOSE_INFO int = 3
const VERBOSE_DEBUG int = 4
const VERBOSE_TRACE int = 5

const DBCREDENTIALS = "dbcredentials"

const CONFIG_FILE = "config"
const CONFIG_FILE_SHORT = "c"

const DOMAIN = "domain"
const DOMAIN_SHORT = "d"

const DOMAIN_FILE string = "domainlist"
const DOMAIN_FILE_SHORT string = "f"


func init() {

	// Set defaults
	//
	// default log loglevel
	viper.SetDefault(VERBOSE, VERBOSE_QUIET)
}
