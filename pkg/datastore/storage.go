package datastore

import (
	"database/sql"

	sl "github.com/eshu0/simplelogger/pkg"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	sl.AppLogger
	database *sql.DB
	Filename string
}

func (fds *Storage) Create() {
	fds.database, _ = sql.Open("sqlite3", fds.Filename)
	CreateHashsTable(fds)
	CreateHashRelationshipTable(fds)

}

func (fds *Storage) Clear() {
	ClearHashData(fds)
	ClearHashRelationshipTable(fds)
}

/*
func db(database *sql.DB) {

	rows, _ := database.Query("SELECT * FROM sqlite_master ") //WHERE type='table' ORDER BY name")
	//var name string

	for rows.Next() {
		//rows.Scan(&name)
		cts, _ := rows.ColumnTypes()
		for _, ct := range cts {
			fmt.Println(ct.Name() + " " + ct.ScanType().Name())
		}

	}

}

*/
