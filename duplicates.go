package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func (fds *DataStorage) GetDuplicateHashes() map[string][]*HashRelationship {
	rows, _ := fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN " + HashsTableName + " on " + HashsTableName + "." + HashsIDColumn + " = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1")

	hashids := fds.ParseDuplicatedHashIDsRows(rows)
	var results map[string][]*HashRelationship
	results = make(map[string][]*HashRelationship)

	for _, hashid := range hashids {
		hrs := fds.GetHashRelationshipByHash(hashid)
		if len(hrs) > 0 {
			hash := hrs[0].Hash.Data
			results[hash] = hrs
		}
	}
	return results
}

func (fds *DataStorage) ParseDuplicatedHashIDsRows(rows *sql.Rows) []int64 {
	var hashid int64
	var count int64

	var results []int64

	for rows.Next() {
		rows.Scan(&hashid, &count)
		results = append(results, hashid)
	}

	return results
}

/*
func (fds *DataStorage) GetDuplicateHashes() {
	rows, _ := fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", " + HashsDataColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN hashes on hashes.id = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1")
}
*/
