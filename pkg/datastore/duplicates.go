package datastore

import (
	"database/sql"
	"fmt"

	"github.com/eshu0/yaft/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

func (fds *Storage) GetDuplicateHashIds(limit int) []models.HashIdnCount {
	var rows *sql.Rows

	if limit < 0 {
		rows, _ = fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN " + HashsTableName + " on " + HashsTableName + "." + HashsIDColumn + " = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1")
	} else {
		rows, _ = fds.database.Query("SELECT " + HashToFilesHashIDColumn + ", COUNT(" + HashToFilesHashIDColumn + ")  FROM " + HashToFilesTableName + " INNER JOIN " + HashsTableName + " on " + HashsTableName + "." + HashsIDColumn + " = " + HashToFilesTableName + "." + HashToFilesHashIDColumn + " GROUP BY " + HashToFilesHashIDColumn + " HAVING COUNT(" + HashToFilesHashIDColumn + ") > 1 LIMIT " + fmt.Sprintf("%d", limit))
	}

	return fds.ParseDuplicatedHashIDsRows(rows)
}

func (fds *Storage) GetDuplicateHashes(limit int) map[string][]*models.HashRelationship {
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

func (fds *Storage) ParseDuplicatedHashIDsRows(rows *sql.Rows) []models.HashIdnCount {
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
