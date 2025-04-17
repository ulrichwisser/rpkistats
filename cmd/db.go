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
	"time"
	"strings"
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

func pushToQueue(db *sql.DB, testid, region, country, sector string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Could not start DB transaction %s", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO QUEUE (testid, region, country, sector) VALUES (?, ?, ?, ?)", testid, region, country, sector)
	if err != nil {
		log.Errorf("Error inserting testid into queue. %s", err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Could not commit to DB %s", err)
	}
}

func popFromQueue(db *sql.DB, testid string) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Could not start DB transaction %s", err)
	}
	defer tx.Rollback()
	_, _ = tx.Exec("DELETE FROM QUEUE WHERE testid = ?", testid)
	log.Debugf("Deleted testid %s from queue.", testid)
	err = tx.Commit()
	if err != nil {
		//log.Fatalf("Could not commit to DB %s", err)
	}
}

func result2int(testresult string) int {
	switch strings.ToUpper(testresult) {
	case "INFO": return 0
	case "NOTICE": return 1
	case "WARNING": return 2
	case "ERROR": return 3
	case "CRITICAL": return 4
	}
	return -1
}

func saveResults(db *sql.DB, testid, region, country, sector string, results map[string]interface{}) error {
	createdAt, _ := time.Parse(time.RFC3339, results["created_at"].(string))
	params := results["params"].(map[string]interface{})
	domain := params["domain"]

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Could not start DB transaction %s", err)
	}
	defer tx.Rollback()

	// Zonemaster results can have more than one entry per testcase
	// we will save it only once, with the worst result of all results for the same testcase
	type Testresult struct {
		Module string
		Testcase string
		Testresult string
	}
	var testresults map[string]Testresult = make(map[string]Testresult, 0)

	for _, entry := range results["results"].([]interface{}) {
		// make our entry easier to work with
		result := entry.(map[string]interface{})
		testcase := result["testcase"].(string)
		testresult := result["level"].(string)

		// if this is a new testcase make an entry in our result map
		if _,ok:=testresults[testcase]; !ok {
			testresults[testcase] = Testresult{Module:result["module"].(string), Testcase: testcase, Testresult: testresult}
		}

		// check if level is worth than what we have
		if result2int(testresult) > result2int(testresults[testcase].Testresult) {
			tt := testresults[testcase]
			tt.Testresult = testresult
			testresults[testcase] = tt
		}
	}

	for testcase := range testresults {
		if result2int(testresults[testcase].Testresult) < result2int("WARNING") {
			continue
		}
		_, err = tx.Exec("INSERT INTO ZONESTATS (testid, domain, region, country, sector, testdate, testmodule, testcase, testresult) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", 
		                  testid, domain, region, country, sector, createdAt, testresults[testcase].Module, testresults[testcase].Testcase, testresults[testcase].Testresult)
		log.Debugf("INSERT INTO ZONESTATS %s, %s, %s, %s, %s, %s, %s, %s, %s", testid, domain, region, country, sector, createdAt, testresults[testcase].Module, testresults[testcase].Testcase, testresults[testcase].Testresult)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Could not commit to DB %s", err)
	}
	return nil
}
