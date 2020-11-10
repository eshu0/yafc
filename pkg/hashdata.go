package yaft

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	sli "github.com/eshu0/simplelogger/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

const HashsTableName = "hashes"
const HashsIDColumn = "id"
const HashsDataColumn = "data"

type HashData struct {
	ID   int64
	Data string
}

func (hd *HashData) String() string {
	return fmt.Sprintf("%d: %s\n", hd.ID, hd.Data)
}

func (hd *HashData) Save(FilePath string, Log sli.ISimpleLogger) bool {
	bytes, err1 := json.MarshalIndent(hd, "", "\t") //json.Marshal(p)
	if err1 != nil {
		Log.LogErrorf("SaveToFile()", "Marshal json for %s failed with %s ", FilePath, err1.Error())
		return false
	}

	err2 := ioutil.WriteFile(FilePath, bytes, 0644)
	if err2 != nil {
		Log.LogErrorf("SaveToFile()", "Saving %s failed with %s ", FilePath, err2.Error())
		return false
	}

	return true
}

func (hd *HashData) CheckFileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false, err
	}
	return !info.IsDir(), nil
}

func (hd *HashData) LoadHashData(FilePath string, Log sli.ISimpleLogger) (*HashData, bool) {
	ok, err := hd.CheckFileExists(FilePath)
	if ok {
		bytes, err1 := ioutil.ReadFile(FilePath) //ReadAll(jsonFile)
		if err1 != nil {
			Log.LogErrorf("LoadFile()", "Reading '%s' failed with %s ", FilePath, err1.Error())
			return nil, false
		}

		vcfs := HashData{}

		err2 := json.Unmarshal(bytes, &vcfs)

		if err2 != nil {
			Log.LogErrorf("LoadFile()", " Loading %s failed with %s ", FilePath, err2.Error())
			return nil, false
		}

		Log.LogDebugf("LoadFile()", "Read ID %d ", vcfs.ID)
		Log.LogDebugf("LoadFile()", "Read Hash %s ", vcfs.Data)

		return &vcfs, true
	} else {

		if err != nil {
			Log.LogErrorf("LoadFile()", "'%s' was not found to load with error: %s", FilePath, err.Error())
		} else {
			Log.LogErrorf("LoadFile()", "'%s' was not found to load", FilePath)
		}

		return nil, false
	}
}

func CreateHashsTable(fds *DataStorage) {
	statement, _ := fds.database.Prepare("CREATE TABLE IF NOT EXISTS " + HashsTableName + " ([" + HashsIDColumn + "] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,[" + HashsDataColumn + "] NVARCHAR(160) NOT NULL UNIQUE)")
	_, err := statement.Exec()
	if err != nil {
		fds.Log.LogErrorE("CreateHashsTable", err)
		return
	}
}

func ClearHashData(fds *DataStorage) {
	statement, _ := fds.database.Prepare("DELETE FROM " + HashsTableName)
	_, err := statement.Exec()
	if err != nil {
		fds.Log.LogErrorE("ClearHashData", err)
		return
	}
}

func (fds *DataStorage) AddHashData(hr *HashData) int64 {
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

func (fds *DataStorage) GetHashData(ID int64) *HashData {
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

func (fds *DataStorage) FindHashData(HashData string) []*HashData {
	statement, _ := fds.database.Prepare("SELECT " + HashsIDColumn + ", " + HashsDataColumn + " FROM " + HashsTableName + " WHERE " + HashsDataColumn + " = ?")
	rows, _ := statement.Query(HashData)
	return fds.ParseHashDataRows(rows)
}

func (fds *DataStorage) GetAllHashData() []*HashData {
	rows, _ := fds.database.Query("SELECT " + HashsIDColumn + ", " + HashsDataColumn + " FROM " + HashsTableName)
	return fds.ParseHashDataRows(rows)
}

func (fds *DataStorage) ParseHashDataRows(rows *sql.Rows) []*HashData {
	var id int64
	var hash string
	results := []*HashData{}

	for rows.Next() {
		rows.Scan(&id, &hash)
		//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + hash)
		fd := HashData{}
		fd.ID = id
		fd.Data = hash
		results = append(results, &fd)
	}
	return results
}
