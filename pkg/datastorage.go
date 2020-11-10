package yaft

import (
	"database/sql"

	sli "github.com/eshu0/simplelogger/pkg/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type DataStorage struct {
	database *sql.DB
	Filename string
	Log      sli.ISimpleLogger
}

func (fds *DataStorage) Create(log sli.ISimpleLogger) {
	fds.database, _ = sql.Open("sqlite3", fds.Filename)
	fds.Log = log
	CreateHashsTable(fds)
	CreateHashRelationshipTable(fds)

}

func (fds *DataStorage) Clear() {
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
