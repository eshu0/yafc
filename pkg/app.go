package yaft

import (
	"flag"

	sl "github.com/eshu0/simplelogger/pkg"
	yaft "github.com/eshu0/yaft/pkg"
)

type YAFTApp struct {
	sl.AppLogger
	fds DataStorage
}

func (fds *YAFTApp) Create() {
	filename := flag.String("logfile", "yaft.log", "Filename out - defaults to yaft.log")
	session := flag.String("sessionid", "123", "Session - defaults to 123")
	dbname := flag.String("db", "./yaft.db", "Database defaults to ./yaft.db")
	inputdir := flag.String("path", "", "")
	cache := flag.String("cache", "", "")
	list := flag.String("list", "", "")
	clear := flag.String("clear", "", "")
	dupes := flag.String("dupes", "", "")
	dupeids := flag.String("dupeids", "", "")
	limit := flag.Int("limit", -1, "")
	savecsv := flag.String("savecsv", "", "")
	hashid := flag.Int("hashid", -1, "")
	filetofind := flag.String("file", "", "File to find")
	deleteifexists := flag.Bool("die", false, "Delete if exists")
	yestoall := flag.Bool("yta", false, "Delete if exists")

	flag.Parse()

	savetocsav := false

	if savecsv != nil && *savecsv != "" {
		savetocsav = true
	}

	slog := sl.NewSimpleLogger(*filename, *session)

	// lets open a flie log using the session
	slog.OpenAllChannels()

	fds := &yaft.DataStorage{}
	fds.Filename = *dbname
	fds.Create(slog)
}
