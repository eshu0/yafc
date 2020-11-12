package models

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	sli "github.com/eshu0/simplelogger/pkg/interfaces"
	_ "github.com/mattn/go-sqlite3"
)

const HashToFilesTableName = "hashtofiles"
const HashToFilesIDColumn = "id"
const HashToFilesHashIDColumn = "hashid"
const HashToFilesPathColumn = "path"
const HashToFilesTypeColumn = "type"

const HashToFilesViewDelete = "type"

//CREATE VIEW IF NOT EXISTS "+HashToFilesViewDelete+" AS "DELETE FROM " + HashToFilesTableName

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
		f.Close()

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
