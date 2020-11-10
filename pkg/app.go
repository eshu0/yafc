package yaft

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	sl "github.com/eshu0/simplelogger/pkg"
)

type YAFTApp struct {
	sl.AppLogger
	FDS       *DataStorage
	Savetocsv bool
	//filename *string
	//session *string
	dbname         *string
	inputdir       *string
	cache          *string
	list           *string
	clear          *string
	dupes          *string
	dupeids        *string
	Limit          *int
	Hashid         *int
	Filetofind     *string
	Deleteifexists *bool
	Yestoall       *bool
}

func (yapp *YAFTApp) ParseFlags() {

	//yapp.filename := flag.String("logfile", "yaft.log", "Filename out - defaults to yaft.log")
	//yapp.session := flag.String("sessionid", "123", "Session - defaults to 123")

	yapp.dbname = flag.String("db", "./yaft.db", "Database defaults to ./yaft.db")
	yapp.inputdir = flag.String("path", "", "")
	yapp.cache = flag.String("cache", "", "")
	yapp.list = flag.String("list", "", "")
	yapp.clear = flag.String("clear", "", "")
	yapp.dupes = flag.String("dupes", "", "")
	yapp.dupeids = flag.String("dupeids", "", "")
	yapp.Limit = flag.Int("limit", -1, "")
	savecsv := flag.String("savecsv", "", "")
	yapp.Hashid = flag.Int("hashid", -1, "")
	yapp.Filetofind = flag.String("file", "", "File to find")
	yapp.Deleteifexists = flag.Bool("die", false, "Delete if exists")
	yapp.Yestoall = flag.Bool("yta", false, "Delete if exists")

	yapp.Savetocsv = false

	if savecsv != nil && *savecsv != "" {
		yapp.Savetocsv = true
	}

	flag.Parse()

}

func (yapp *YAFTApp) Create() {
	yapp.FDS = &DataStorage{}
	yapp.FDS.Filename = *yapp.dbname
	yapp.FDS.Create(yapp.Log)
}

func (yapp *YAFTApp) Process() {

	if yapp.Filetofind != nil && *yapp.Filetofind != "" {
		reader := bufio.NewReader(os.Stdin)

		if yapp.Deleteifexists != nil && *yapp.Deleteifexists {
			fmt.Println("Delete if exists")
		}

		if yapp.Yestoall != nil && *yapp.Yestoall {
			fmt.Println("Yes to all")
		}

		err := filepath.Walk(*yapp.Filetofind, CompareDirectory(yapp.FDS, yapp.Log, yapp.Deleteifexists, yapp.Yestoall, reader))
		if err != nil {
			panic(err)
		}
	}

	if yapp.Inputdir != nil && *yapp.Inputdir != "" {

		persist := (yapp.cache != nil && *yapp.cache != "")

		err := filepath.Walk(*yapp.inputdir, WalkDir(yapp.FDS, yapp.Log, persist))
		if err != nil {
			panic(err)
		}
	}

	if yapp.list != nil && *yapp.list != "" {
		fmt.Println("Listing all ")
		results := yapp.FDS.GetAllHashData()

		fmt.Println("Hashs: ")
		for _, hd := range results {
			fmt.Println(hd)
		}

		fmt.Println("Relations: ")
		results2 := yapp.FDS.GetAllHashRelationships()
		for _, hr := range results2 {
			fmt.Println(hr)
		}

	}

	if yapp.clear != nil && *yapp.clear != "" {
		yapp.FDS.Clear()
	}

	if yapp.dupes != nil && *yapp.dupes != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if yapp.savetocsav {
			file, err = os.Create("results.csv")
			if err != nil {
				yapp.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
				return
			}
			defer file.Close()
		}

		fmt.Println("Duplicates: ")
		limitcount := -1
		if yapp.imit != nil && *yapp.limit > 0 {
			limitcount = *yapp.limit
		}

		results1 := yapp.FDS.GetDuplicateHashes(limitcount)

		SaveDuplicates("./results.json", yapp.Log, results1)
		if yapp.Savetocsav {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		for k, v := range results1 {
			fmt.Println("Key ", k)
			for _, hr := range v {
				fmt.Println(hr)

				if yapp.Savetocsav {
					err := writer.Write(hr.CSV())
					if err != nil {
						yapp.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
						break
					}
				}

			}
		}
	}

	if dupeids != nil && *dupeids != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if savetocsav {
			file, err = os.Create("ids.csv")
			if err != nil {
				slog.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
				return
			}

			defer file.Close()
		}

		fmt.Println("Duplicate Ids: ")
		limitcount := -1
		if limit != nil && *limit > 0 {
			limitcount = *limit
		}

		results1 := fds.GetDuplicateHashIds(limitcount)

		if savetocsav {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		var res []string
		if savetocsav {

			res = []string{}
			//
			res = append(res, "Hash Id")
			res = append(res, "Count")

			err := writer.Write(res)
			if err != nil {
				slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
				return
			}
		}

		for k, v := range results1 {
			fmt.Printf("%d id = %d \n", k, v)
			if savetocsav {

				res = []string{}
				//res = append(res,	fmt.Sprintf("%d", k))
				res = append(res, fmt.Sprintf("%d", v.HashId))
				res = append(res, fmt.Sprintf("%d", v.Count))

				err := writer.Write(res)
				if err != nil {
					yapp.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
					break
				}
			}

		}

	}

	if hashid != nil && *hashid > 0 {

		var file *os.File
		var writer *csv.Writer
		var err error

		if savetocsav {
			file, err = os.Create("files.csv")
			if err != nil {
				yapp.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
				return
			}

			defer file.Close()
		}

		fmt.Println("Files: ")

		results1 := yapp.FDS.GetFilesByHashId(int64(*yapp.hashid))

		if yapp.savetocsav {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		var res []string
		if yapp.savetocsav {

			res = []string{}
			//
			res = append(res, "Hash Id")
			res = append(res, "Files")

			err := writer.Write(res)
			if err != nil {
				yapp.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
				return
			}
		}

		for k, v := range results1 {
			fmt.Printf("%d) %d = %s \n", k, *yapp.hashid, v)
			if yapp.savetocsav {

				res = []string{}
				//res = append(res,	fmt.Sprintf("%d", k))
				res = append(res, fmt.Sprintf("%d", *yapp.hashid))
				res = append(res, v)

				err := writer.Write(res)
				if err != nil {
					yapp.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
					break
				}
			}

		}

	}
}
