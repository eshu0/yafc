package datastore

import (
	"database/sql"

	"github.com/eshu0/yaft/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

const HashsTableName = "hashes"
const HashsIDColumn = "id"
const HashsDataColumn = "data"

func CreateHashsTable(fds *Storage) {
	statement, _ := fds.database.Prepare("CREATE TABLE IF NOT EXISTS " + HashsTableName + " ([" + HashsIDColumn + "] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,[" + HashsDataColumn + "] NVARCHAR(160) NOT NULL UNIQUE)")
	_, err := statement.Exec()
	if err != nil {
		fds.Log.LogErrorE("CreateHashsTable", err)
		return
	}
}

func ClearHashData(fds *Storage) {
	statement, _ := fds.database.Prepare("DELETE FROM " + HashsTableName)
	_, err := statement.Exec()
	if err != nil {
		fds.Log.LogErrorE("ClearHashData", err)
		return
	}
}

func (fds *Storage) AddHashData(hr *models.HashData) int64 {
	statement, _ := fds.database.Prepare("INSERT INTO " + HashsTableName + " (" + HashsDataColumn + ") VALUES (?)")
	res, err := statement.Exec(hr.Data)
	if err == nil {
		lastid, _ := res.LastInsertId()
		return lastid
	} else {
		fds.Log.LogErrorE("AddHashData", err)
		return -1
	}

}

func (fds *Storage) GetHashData(ID int64) *models.HashData {
	statement, _ := fds.database.Prepare("SELECT " + HashsIDColumn + ", " + HashsDataColumn + " FROM " + HashsTableName + " WHERE " + HashsIDColumn + " = ?")
	rows, err := statement.Query(ID)
	if err != nil {
		fds.Log.LogErrorE("GetHashData", err)
		return nil
	}
	res := fds.ParseHashDataRows(rows)
	if len(res) == 0 {
		return nil
	} else {
		return res[0]
	}
}

func (fds *Storage) FindHashData(HashData string) []*models.HashData {
	statement, _ := fds.database.Prepare("SELECT " + HashsIDColumn + ", " + HashsDataColumn + " FROM " + HashsTableName + " WHERE " + HashsDataColumn + " = ?")
	rows, _ := statement.Query(HashData)
	return fds.ParseHashDataRows(rows)
}

func (fds *Storage) GetAllHashData() []*models.HashData {
	rows, _ := fds.database.Query("SELECT " + HashsIDColumn + ", " + HashsDataColumn + " FROM " + HashsTableName)
	return fds.ParseHashDataRows(rows)
}

func (fds *Storage) ParseHashDataRows(rows *sql.Rows) []*models.HashData {
	var id int64
	var hash string
	results := []*models.HashData{}

	for rows.Next() {
		rows.Scan(&id, &hash)
		//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + hash)
		fd := models.HashData{}
		fd.ID = id
		fd.Data = hash
		results = append(results, &fd)
	}
	return results
}
