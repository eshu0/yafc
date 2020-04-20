package main

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"os"

	sli "github.com/eshu0/simplelogger/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

const HashToFilesTableName = "hashtofiles"
const HashToFilesIDColumn = "id"
const HashToFilesHashIDColumn = "hashid"
const HashToFilesPathColumn = "path"
const HashToFilesTypeColumn = "type"

type HashRelationship struct {
	ID   int64
	Path string
	Type int
	Hash *HashData
}

func (hr *HashRelationship) GenHashData(Logger sli.ISimpleLogger, FilePath string, isdir bool) bool {

	f, err := os.Open(FilePath)

	if err != nil {
		Logger.LogErrorE("Visit", err)
		return false
	} else {

		input := bufio.NewReader(f)

		hash := sha256.New()
		if _, err := io.Copy(hash, input); err != nil {
			Logger.LogErrorE("Visit", err)
		}
		sum := hash.Sum(nil)

		//fmt.Printf("%s %x\n", path, sum)
		Logger.LogInfo("Visit", fmt.Sprintf("%s %x", FilePath, sum))

		if hr.Hash == nil {
			fh := HashData{}
			fh.ID = -1
			fh.Data = fmt.Sprintf("%x", sum)
			hr.Hash = &fh
		}

		//hr.addFilepath(Logger, FilePath, isdir)
		if isdir {
			hr.Type = 0
		} else {
			hr.Type = 1
		}

		return true
	}

}

/*
func (hr *HashRelationship) addFilepath(Logger sli.ISimpleLogger, path string, isdir bool) {

	files := hr.FilePaths
	fn := FileData{}
	fn.ID = -1
	fn.Path = path



	files = append(files, &fn)
	hr.FilePaths = files

}
*/
func CreateHashRelationshipTable(fds *DataStorage) {
	statement, _ := fds.database.Prepare("CREATE TABLE IF NOT EXISTS " + HashToFilesTableName + " ([" + HashToFilesIDColumn + "] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, [" + HashToFilesHashIDColumn + "] INTEGER REFERENCES " + HashsTableName + "(" + HashsIDColumn + ") ,[" + HashToFilesPathColumn + "] TEXT NOT NULL UNIQUE, [" + HashToFilesTypeColumn + "] INTEGER NOT NULL)")
	statement.Exec()
}

func ClearHashRelationshipTable(fds *DataStorage) {
	statement, _ := fds.database.Prepare("DELETE FROM " + HashToFilesTableName)
	statement.Exec()
}

func (fds *DataStorage) AddHashRelationship(hr *HashRelationship) { //(int64, []int64) {

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

func (fds *DataStorage) GetAllHashRelationships() []*HashRelationship {
	rows, _ := fds.database.Query("SELECT " + HashToFilesIDColumn + ", " + HashToFilesHashIDColumn + ", " + HashToFilesPathColumn + ", " + HashToFilesTypeColumn + " FROM " + HashToFilesTableName)
	return fds.ParseHashRelationshipRows(rows)
}

func (fds *DataStorage) GetHashRelationshipByHash(hashid int64) []*HashRelationship {
	statement, _ := fds.database.Prepare("SELECT " + HashToFilesIDColumn + ", " + HashToFilesHashIDColumn + ", " + HashToFilesPathColumn + ", " + HashToFilesTypeColumn + " FROM " + HashToFilesTableName + " WHERE " + HashToFilesHashIDColumn + " = ? ")
	rows, _ := statement.Query(hashid)
	return fds.ParseHashRelationshipRows(rows)
}

func (fds *DataStorage) ParseHashRelationshipRows(rows *sql.Rows) []*HashRelationship {
	var id int64
	var hashid int64
	var path string
	var typei int

	var results []*HashRelationship
	var lasthash *HashData

	for rows.Next() {
		rows.Scan(&id, &hashid, &path, &typei)
		hr := HashRelationship{}
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

func (fds *DataStorage) ParseHashRelationshipRows1(rows *sql.Rows) map[string]*HashRelationship {
	var id int64
	var hashid int64
	var path string
	var typei int

	var results map[string]*HashRelationship
	results = make(map[string]*HashRelationship)

	var lasthash *HashData

	for rows.Next() {
		rows.Scan(&id, &hashid, &path, &typei)

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
