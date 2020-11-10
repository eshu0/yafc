package yaft

import (
	"flag"

	sl "github.com/eshu0/simplelogger/pkg"
)

type YAFTApp struct {
	sl.AppLogger
	FDS *DataStorage
	Savetocsv bool
	//filename *string
	//session *string
	dbname *string
	inputdir *string
	cache *string
	list *string
	clear *string
	dupes *string
	dupeids *string
	Limit *int
	Hashid *int
	Filetofind *string
	Deleteifexists *bool
	Yestoall bool	
}

func (yapp *YAFTApp) ParseFlags() {
	
	//yapp.filename := flag.String("logfile", "yaft.log", "Filename out - defaults to yaft.log")
	//yapp.session := flag.String("sessionid", "123", "Session - defaults to 123")

	yapp.dbname := flag.String("db", "./yaft.db", "Database defaults to ./yaft.db")
	yapp.nputdir := flag.String("path", "", "")
	yapp.cache := flag.String("cache", "", "")
	yapp.list := flag.String("list", "", "")
	yapp.clear := flag.String("clear", "", "")
	yapp.dupes := flag.String("dupes", "", "")
	yapp.dupeids := flag.String("dupeids", "", "")
	yapp.Limit := flag.Int("limit", -1, "")
	savecsv := flag.String("savecsv", "", "")
	yapp.Hashid := flag.Int("hashid", -1, "")
	yapp.Filetofind := flag.String("file", "", "File to find")
	yapp.Deleteifexists := flag.Bool("die", false, "Delete if exists")
	yestoall := flag.Bool("yta", false, "Delete if exists")

	yapp.Savetocsv = false

	if savecsv != nil && *savecsv != "" {
		yapp.Savecsv = true
	}

	
	yapp.Yestoall = false
	if yestoall != nil {
		yapp.Yestoall = *yestoall
	}


	flag.Parse()

}

func (yapp *YAFTApp) Create() {


	yapp.FDS := &DataStorage{}
	yapp.FDS.Filename = *yapp.dbname
	yapp.FDS.Create(slog)
}
