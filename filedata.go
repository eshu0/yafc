package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const FilesTableName = "files"
const FilesIDColumn = "id"
const FilesPathColumn = "path"
const FilesTypeColumn = "type"

type FileData struct {
	ID   int64
	Path string
	Type int
}

func (fd *FileData) String() string {
	return fmt.Sprintf("%d (%d): %s\n", fd.ID, fd.Type, fd.Path)
}

func CreateFilesTable(fds *DataStorage) {
	statement, _ := fds.database.Prepare("CREATE TABLE IF NOT EXISTS " + FilesTableName + " ([" + FilesIDColumn + "] INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,[" + FilesPathColumn + "] TEXT NOT NULL, [" + FilesTypeColumn + "] INTEGER NOT NULL)")
	statement.Exec()
}

func (fds *DataStorage) AddFileData(hr *FileData) int64 {
	statement, _ := fds.database.Prepare("INSERT INTO " + FilesTableName + " (" + FilesPathColumn + " , " + FilesTypeColumn + ") VALUES (?, ?)")
	res, _ := statement.Exec(hr.Path, hr.Type)
	lastid, _ := res.LastInsertId()
	return lastid
}

func (fds *DataStorage) GetFileData(ID int64) *FileData {
	statement, _ := fds.database.Prepare("SELECT " + FilesIDColumn + ", " + FilesPathColumn + ", " + FilesTypeColumn + " FROM " + FilesTableName + " WHERE " + FilesIDColumn + " = ?")
	rows, _ := statement.Query(ID)
	res := fds.ParseFileDataRows(rows)
	if len(res) == 0 {
		return nil
	} else {
		return res[0]
	}
}

func (fds *DataStorage) FindFileData(Path string) []*FileData {
	statement, _ := fds.database.Prepare("SELECT " + FilesIDColumn + ", " + FilesPathColumn + ", " + FilesTypeColumn + " FROM " + FilesTableName + " WHERE " + FilesPathColumn + " = ?")
	rows, _ := statement.Query(Path)
	return fds.ParseFileDataRows(rows)
}

func (fds *DataStorage) GetAllFileData() []*FileData {
	rows, _ := fds.database.Query("SELECT " + FilesIDColumn + ", " + FilesPathColumn + ", " + FilesTypeColumn + " FROM " + FilesTableName)
	return fds.ParseFileDataRows(rows)
}

func (fds *DataStorage) ParseFileDataRows(rows *sql.Rows) []*FileData {
	var id int64
	var path string
	var typei int
	results := []*FileData{}

	for rows.Next() {
		rows.Scan(&id, &path, &typei)
		//	fmt.Println("READ: " + strconv.Itoa(id) + ": " + path )
		fd := FileData{}
		fd.ID = id
		fd.Path = path
		fd.Type = typei
		results = append(results, &fd)
	}
	return results
}
