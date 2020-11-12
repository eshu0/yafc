package models

import (
	"encoding/json"
	"io/ioutil"

	sli "github.com/eshu0/simplelogger/pkg/interfaces"
	"github.com/eshu0/yaft/pkg/models"
	_ "github.com/mattn/go-sqlite3"
)

type HashIdnCount struct {
	HashId int64
	Count  int64
}

func SaveDuplicates(FilePath string, Log sli.ISimpleLogger, hd map[string][]*models.HashRelationship) bool {
	bytes, err1 := json.MarshalIndent(hd, "", "\t") //json.Marshal(p)
	if err1 != nil {
		Log.LogErrorf("SaveDuplicates()", "Marshal json for %s failed with %s ", FilePath, err1.Error())
		return false
	}

	err2 := ioutil.WriteFile(FilePath, bytes, 0644)
	if err2 != nil {
		Log.LogErrorf("SaveDuplicates()", "Saving %s failed with %s ", FilePath, err2.Error())
		return false
	}

	return true
}
