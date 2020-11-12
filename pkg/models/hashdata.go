package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	sli "github.com/eshu0/simplelogger/pkg/interfaces"
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
		Log.LogErrorf("Save()", "Marshal json for %s failed with %s ", FilePath, err1.Error())
		return false
	}

	err2 := ioutil.WriteFile(FilePath, bytes, 0644)
	if err2 != nil {
		Log.LogErrorf("Save()", "Saving %s failed with %s ", FilePath, err2.Error())
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
			Log.LogErrorf("LoadHashData()", "Reading '%s' failed with %s ", FilePath, err1.Error())
			return nil, false
		}

		vcfs := HashData{}

		err2 := json.Unmarshal(bytes, &vcfs)

		if err2 != nil {
			Log.LogErrorf("LoadHashData()", " Loading %s failed with %s ", FilePath, err2.Error())
			return nil, false
		}

		Log.LogDebugf("LoadHashData()", "Read ID %d ", vcfs.ID)
		Log.LogDebugf("LoadHashData()", "Read Hash %s ", vcfs.Data)

		return &vcfs, true
	} else {

		if err != nil {
			Log.LogErrorf("LoadHashData()", "'%s' was not found to load with error: %s", FilePath, err.Error())
		} else {
			Log.LogErrorf("LoadHashData()", "'%s' was not found to load", FilePath)
		}

		return nil, false
	}
}
