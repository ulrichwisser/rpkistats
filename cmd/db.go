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
	"github.com/spf13/viper"
	"github.com/apex/log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func openDB() *sql.DB {
	// open database
	if viper.GetString(DBCREDENTIALS) == "" {
		log.Fatal("No DB credentials given.")
	}
	db, err := sql.Open("mysql", viper.GetString(DBCREDENTIALS))
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Debug("DB OPEN")
	return db
}

func rpki2db(db *sql.DB, rpkistats []*RPKIstat) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Could not start DB transaction %s", err)
	}
	defer tx.Rollback()

	for _,rpki := range rpkistats {
		_, err = tx.Exec("INSERT INTO RPKI(TESTDATE,TLD,NAMES,NAMES_ROA_FULL,NAMES_ROA_PARTIAL,IP4S,IP4S_ROAS,IP6s, IP6S_ROAS,TAS4, TAS6, AS4, AS6) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", 
		                  rpki.Date, rpki.Domain, rpki.Names, rpki.NamesFull, rpki.NamesPartial, rpki.IPv4, rpki.IPv4roas, rpki.IPv6, rpki.IPv6roas, rpki.TAs4, rpki.TAs6, rpki.AS4, rpki.AS6)
		log.Debugf("INSERT INTO RPKI %s, %-15s, Names %2d, Full %2d, Partail %2d, IPv4 %2d, ROA %2d, IPv6 %2d, ROA %2d, TA4 %1d, TA6 %1d, AS4 %2d, AS6 %2d", 
		                  rpki.Date, rpki.Domain, rpki.Names, rpki.NamesFull, rpki.NamesPartial, rpki.IPv4, rpki.IPv4roas, rpki.IPv6, rpki.IPv6roas, rpki.TAs4, rpki.TAs6, rpki.AS4, rpki.AS6)
		if err != nil {
			log.Fatalf("Could not insert into %s", err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Could not commit to DB %s", err)
	} 
	log.Debug("Data committed to database")
	return
}
