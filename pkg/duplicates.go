package yaft

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"

	sli "github.com/eshu0/simplelogger/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

type HashIdnCount struct {
	HashId int64
	Count  int64
}

func SaveDuplicates(FilePath string, Log sli.ISimpleLogger, hd map[string][]*HashRelationship) bool {
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

func (fds *DataStorage) GetDuplicateHashIds(limit int) []HashIdnCount {
	var rows *sql.Rows

	if limit < 0 {
		rows, _ = fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN " + HashsTableName + " on " + HashsTableName + "." + HashsIDColumn + " = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1")
	} else {
		rows, _ = fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN " + HashsTableName + " on " + HashsTableName + "." + HashsIDColumn + " = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1 LIMIT " + fmt.Sprintf("%d", limit))
	}

	return fds.ParseDuplicatedHashIDsRows(rows)
}

func (fds *DataStorage) GetDuplicateHashes(limit int) map[string][]*HashRelationship {
	ids := fds.GetDuplicateHashIds(limit)
	var results map[string][]*HashRelationship
	results = make(map[string][]*HashRelationship)

	for _, id := range ids {
		hrs := fds.GetHashRelationshipByHash(id.HashId)
		if len(hrs) > 0 {
			hash := hrs[0].Hash.Data
			results[hash] = hrs
		}
	}
	return results
}

func (fds *DataStorage) ParseDuplicatedHashIDsRows(rows *sql.Rows) []HashIdnCount {
	var hashid int64
	var count int64

	var results []HashIdnCount

	for rows.Next() {
		rows.Scan(&hashid, &count)
		hic := HashIdnCount{}
		hic.Count = count
		hic.HashId = hashid
		results = append(results, hic)
	}

	return results
}

/*
func (fds *DataStorage) GetDuplicateHashes() {
	rows, _ := fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", " + HashsDataColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN hashes on hashes.id = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1")
}
*/
