package datastore

import (
	"database/sql"
	"fmt"

	"github.com/eshu0/yaft/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

const HashToFilesTableName = "hashtofiles"
const HashToFilesIDColumn = "id"
const HashToFilesHashIDColumn = "hashid"
const HashToFilesPathColumn = "path"
const HashToFilesTypeColumn = "type"

const HashToFilesViewDelete = "type"

//CREATE VIEW IF NOT EXISTS "+HashToFilesViewDelete+" AS "DELETE FROM " + HashToFilesTableName

func CreateHashRelationshipTable(fds *Storage) {
	statement, _ := fds.database.Prepare("CREATE TABLE IF NOT EXISTS " + HashToFilesTableName + " ([" + HashToFilesIDColumn + "] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, [" + HashToFilesHashIDColumn + "] INTEGER REFERENCES " + HashsTableName + "(" + HashsIDColumn + ") ,[" + HashToFilesPathColumn + "] TEXT NOT NULL UNIQUE, [" + HashToFilesTypeColumn + "] INTEGER NOT NULL)")
	statement.Exec()
}

func ClearHashRelationshipTable(fds *Storage) {
	statement, _ := fds.database.Prepare("DELETE FROM " + HashToFilesTableName)
	statement.Exec()
}

func (fds *Storage) AddHashRelationship(hr *models.HashRelationship) { //(int64, []int64) {

	var hidtoinsert int64
	//var fidtoinsert int64

	results := fds.FindHashData(hr.Hash.Data)
	if len(results) == 0 {
		hidtoinsert = fds.AddHashData(hr.Hash)
	} else {
		hidtoinsert = results[0].ID
	}
	hr.Hash.ID = hidtoinsert
	/*
		for _, fd := range hr.FilePaths {
			results := fds.FindFileData(fd.Path)
			if len(results) == 0 {
				fidtoinsert = fds.AddFileData(fd)
			} else {
				fidtoinsert = results[0].ID
			}
	*/
	//hrres := fds.GetHashRelationshipByHash(hidtoinsert)
	//if len(hrres) == 0 {
	//fmt.Println(fmt.Sprintf("adding %d to %d", hidtoinsert, fidtoinsert))
	fmt.Println(fmt.Sprintf("adding %d to %s", hidtoinsert, hr.Path))

	//statement, _ := fds.database.Prepare("INSERT INTO " + HashToFilesTableName + " (" + HashToFilesHashIDColumn + "," + HashToFilesFilesIDColumn + ") VALUES (?,?)")
	//statement.Exec(hidtoinsert, fidtoinsert)
	statement, _ := fds.database.Prepare("INSERT INTO " + HashToFilesTableName + " (" + HashToFilesHashIDColumn + "," + HashToFilesPathColumn + "," + HashToFilesTypeColumn + ") VALUES (?,?,?)")
	statement.Exec(hidtoinsert, hr.Path, hr.Type)
	//} else {
	//	for _, hre := range hrres {
	//		fmt.Println(fmt.Sprintf("Was adding %d to %d Relationship however it exists at %s", hidtoinsert, fidtoinsert, hre))
	//	}
	//}

	//}

}

func (fds *Storage) GetFilesByHashId(hashid int64) []string {
	statement, _ := fds.database.Prepare("SELECT " + HashToFilesPathColumn + " FROM " + HashToFilesTableName + " WHERE " + HashToFilesHashIDColumn + " = ? ")
	rows, _ := statement.Query(hashid)
	var path string
	var results []string

	for rows.Next() {
		rows.Scan(&path)
		results = append(results, path)
	}

	return results
}

func (fds *Storage) GetAllHashRelationships() []*models.HashRelationship {
	rows, _ := fds.database.Query("SELECT " + HashToFilesIDColumn + ", " + HashToFilesHashIDColumn + ", " + HashToFilesPathColumn + ", " + HashToFilesTypeColumn + " FROM " + HashToFilesTableName)
	return fds.ParseHashRelationshipRows(rows)
}

func (fds *Storage) GetHashRelationshipByHash(hashid int64) []*models.HashRelationship {
	statement, _ := fds.database.Prepare("SELECT " + HashToFilesIDColumn + ", " + HashToFilesHashIDColumn + ", " + HashToFilesPathColumn + ", " + HashToFilesTypeColumn + " FROM " + HashToFilesTableName + " WHERE " + HashToFilesHashIDColumn + " = ? ")
	rows, _ := statement.Query(hashid)
	return fds.ParseHashRelationshipRows(rows)
}

func (fds *Storage) ParseHashRelationshipRows(rows *sql.Rows) []*models.HashRelationship {
	var id int64
	var hashid int64
	var path string
	var typei int

	var results []*models.HashRelationship
	var lasthash *models.HashData

	for rows.Next() {
		rows.Scan(&id, &hashid, &path, &typei)
		hr := models.HashRelationship{}
		hr.ID = id
		hr.Path = path
		hr.Type = typei

		if lasthash == nil {
			//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + hash)
			lasthash = fds.GetHashData(hashid)
			if lasthash == nil {
				break
			}
			/*
				if results[lasthash.Data] == nil {
					hr.Hash = lasthash
					results[lasthash.Data] = &hr
				}
			*/
		} else {

			if lasthash.ID != hashid {
				lasthash = fds.GetHashData(hashid)
				/*
					if results[lasthash.Data] == nil {
						hr := HashRelationship{}
						hr.Hash = lasthash
						results[lasthash.Data] = &hr
					}
				*/
			}
		}
		hr.Hash = lasthash
		results = append(results, &hr)
	}

	return results
}

func (fds *Storage) ParseHashRelationshipRows1(rows *sql.Rows) map[string]*models.HashRelationship {
	var id int64
	var hashid int64
	var path string
	var typei int

	var results map[string]*models.HashRelationship
	results = make(map[string]*models.HashRelationship)

	var lasthash *models.HashData

	for rows.Next() {
		rows.Scan(&id, &hashid, &path, &typei)

		if lasthash == nil {
			//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + hash)
			lasthash = fds.GetHashData(hashid)
			if lasthash == nil {
				break
			}

			if results[lasthash.Data] == nil {
				hr := models.HashRelationship{}
				hr.Hash = lasthash
				results[lasthash.Data] = &hr
			}

		} else {
			if lasthash.ID != hashid {
				lasthash = fds.GetHashData(hashid)

				if results[lasthash.Data] == nil {
					hr := models.HashRelationship{}
					hr.Hash = lasthash
					results[lasthash.Data] = &hr
				}
			}
		}
	}
	return results
}

/*
func (fds *DataStorage) ParseHashRelationshipRows(rows *sql.Rows) map[string]*HashRelationship {
	var id int64
	var hashid int64
	var fileid int64

	var results map[string]*HashRelationship
	results = make(map[string]*HashRelationship)

	var lasthash *HashData

	for rows.Next() {
		rows.Scan(&id, &hashid, &fileid)

		if lasthash == nil {
			//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + hash)
			lasthash = fds.GetHashData(hashid)
			if lasthash == nil {
				break
			}

			if results[lasthash.Data] == nil {
				hr := HashRelationship{}
				hr.Hash = lasthash
				results[lasthash.Data] = &hr
			}

		} else {
			if lasthash.ID != hashid {
				lasthash = fds.GetHashData(hashid)

				if results[lasthash.Data] == nil {
					hr := HashRelationship{}
					hr.Hash = lasthash
					results[lasthash.Data] = &hr
				}
			}
		}

		hr := results[lasthash.Data]
		filedata := fds.GetFileData(fileid)

		if filedata == nil {

			foundfiledata := false
			for _, fd := range hr.FilePaths {
				if fd.ID == filedata.ID {
					foundfiledata = true
					break
				}
			}

			if !foundfiledata {
				files := hr.FilePaths
				files = append(files, filedata)
				hr.FilePaths = files
			}
		} else {
			files := hr.FilePaths
			files = append(files, filedata)
			hr.FilePaths = files
		}

		ids := hr.IDs
		ids = append(ids, id)
		hr.IDs = ids
	}
	return results
}
*/
func (hr *HashRelationship) String() string {
	b := ":::: ID ::: \n"
	b += fmt.Sprintf("\t %d\n", hr.ID)
	b += ":::: Hash ::: \n"
	b += fmt.Sprintf("\t %s", hr.Hash)
	b += ":::: File ::: \n"
	b += fmt.Sprintf("\t %s\n", hr.Path)

	return b
}

func (hr *HashRelationship) CSV() []string {
	b := []string{}
	b = append(b, fmt.Sprintf("%d", hr.ID))
	b = append(b, fmt.Sprintf("%s", hr.Hash))
	b = append(b, fmt.Sprintf("%s", hr.Path))
	return b
}
