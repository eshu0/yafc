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

	FDS            *DataStorage
	Savetocsv      bool
	DBFileName     *string
	Inputdir       *string
	Cache          *string
	List           *string
	Clear          *string
	Dupes          *string
	Dupeids        *string
	Limit          *int
	Hashid         *int
	Filetofind     *string
	Deleteifexists *bool
	Yestoall       *bool
}

func (yapp *YAFTApp) ParseFlags() {

	yapp.LogInfo("ParseFlags Started")

	yapp.DBFileName = flag.String("db", "./yaft.db", "Database defaults to ./yaft.db")
	yapp.Inputdir = flag.String("path", "", "")
	yapp.Cache = flag.String("cache", "", "")
	yapp.List = flag.String("list", "", "")
	yapp.Clear = flag.String("clear", "", "")
	yapp.Dupes = flag.String("dupes", "", "")
	yapp.Dupeids = flag.String("dupeids", "", "")

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

	yapp.LogInfo("ParseFlags Finished")
}

func (yapp *YAFTApp) Create() {
	yapp.LogInfo("Create", "Creating Database Started")

	yapp.FDS = &DataStorage{}
	yapp.FDS.Filename = *yapp.DBFileName
	yapp.FDS.Create()

	yapp.LogInfo("Create", "Creating Database Finished")
}

func (yapp *YAFTApp) Process() {
	yapp.LogInfo("Process", "Process Started")

	if yapp.Filetofind != nil && *yapp.Filetofind != "" {
		reader := bufio.NewReader(os.Stdin)

		if yapp.Deleteifexists != nil && *yapp.Deleteifexists {
			yapp.LogInfo("Filetofind", "Delete if exists")
		}

		if yapp.Yestoall != nil && *yapp.Yestoall {
			yapp.LogInfo("Filetofind", "Yes to all")
		}

		err := filepath.Walk(*yapp.Filetofind, CompareDirectory(yapp.FDS, yapp.Log, yapp.Deleteifexists, yapp.Yestoall, reader))
		if err != nil {
			panic(err)
		}
	}

	if yapp.Inputdir != nil && *yapp.Inputdir != "" {
		yapp.LogInfo("Inputdir", "Input dir")
		persist := (yapp.Cache != nil && *yapp.Cache != "")

		err := filepath.Walk(*yapp.Inputdir, WalkDir(yapp.FDS, yapp.Log, persist))
		if err != nil {
			panic(err)
		}
	}

	if yapp.List != nil && *yapp.List != "" {
		yapp.LogInfo("List", "Listing all ")
		results := yapp.FDS.GetAllHashData()

		yapp.LogInfo("List", "Hashs: ")
		for _, hd := range results {
			yapp.LogInfo("List", hd)
		}

		yapp.LogInfo("List", "Relations: ")
		results2 := yapp.FDS.GetAllHashRelationships()
		for _, hr := range results2 {
			yapp.LogInfo("List", hr)
		}

	}

	if yapp.Clear != nil && *yapp.Clear != "" {
		yapp.LogInfo("Clear", "Clearing database")
		yapp.FDS.Clear()
	}

	if yapp.Dupes != nil && *yapp.Dupes != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if yapp.Savetocsv {
			file, err = os.Create("results.csv")
			if err != nil {
				yapp.LogError("Duplicates", "Cannot create file")
				yapp.LogErrorE("Duplicates", err)
				return
			}
			defer file.Close()
		}

		yapp.LogInfo("Duplicates", "Duplicates: ")
		limitcount := -1
		if yapp.Limit != nil && *yapp.Limit > 0 {
			limitcount = *yapp.Limit
		}

		results1 := yapp.FDS.GetDuplicateHashes(limitcount)

		SaveDuplicates("./results.json", yapp.Log, results1)
		if yapp.Savetocsv {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		for k, v := range results1 {
			yapp.LogInfo("Duplicates", "Key ", k)
			for _, hr := range v {
				yapp.LogInfo("Duplicates", hr)

				if yapp.Savetocsv {
					err := writer.Write(hr.CSV())
					if err != nil {
						yapp.LogError("Duplicates", "Cannot save to CSV")
						yapp.LogErrorE("Duplicates", err)
						break
					}
				}

			}
		}
	}

	if yapp.Dupeids != nil && *yapp.Dupeids != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if yapp.Savetocsv {
			file, err = os.Create("ids.csv")
			if err != nil {
				yapp.LogError("DuplicateIds", "Cannot save to ids.csv")
				yapp.LogErrorE("DuplicateIds", err)
				return
			}

			defer file.Close()
		}

		yapp.LogInfo("DuplicateIds", "Duplicate Ids: ")
		limitcount := -1
		if yapp.Limit != nil && *yapp.Limit > 0 {
			limitcount = *yapp.Limit
		}

		results1 := yapp.FDS.GetDuplicateHashIds(limitcount)

		if yapp.Savetocsv {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		var res []string
		if yapp.Savetocsv {

			res = []string{}
			//
			res = append(res, "Hash Id")
			res = append(res, "Count")

			err := writer.Write(res)
			if err != nil {
				yapp.LogError("DuplicateIds", "Cannot save to ids.csv")
				yapp.LogErrorE("DuplicateIds", err)
				return
			}
		}

		for k, v := range results1 {
			yapp.LogInfof("DuplicateIds", "%d id = %d \n", k, v)
			if yapp.Savetocsv {

				res = []string{}
				//res = append(res,	fmt.Sprintf("%d", k))
				res = append(res, fmt.Sprintf("%d", v.HashId))
				res = append(res, fmt.Sprintf("%d", v.Count))

				err := writer.Write(res)
				if err != nil {
					yapp.LogError("DuplicateIds", "Cannot save to ids.csv")
					yapp.LogErrorE("DuplicateIds", err)
					break
				}
			}

		}

	}

	if yapp.Hashid != nil && *yapp.Hashid > 0 {

		var file *os.File
		var writer *csv.Writer
		var err error

		if yapp.Savetocsv {
			file, err = os.Create("files.csv")
			if err != nil {
				yapp.LogError("Hashid", "Cannot save to files.csv")
				yapp.LogErrorE("Hashid", err)
				return
			}

			defer file.Close()
		}

		yapp.LogInfo("Hashid", "Files: ")

		results1 := yapp.FDS.GetFilesByHashId(int64(*yapp.Hashid))

		if yapp.Savetocsv {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		var res []string
		if yapp.Savetocsv {

			res = []string{}
			//
			res = append(res, "Hash Id")
			res = append(res, "Files")

			err := writer.Write(res)
			if err != nil {
				yapp.LogError("Hashid", "Cannot save to files.csv")
				yapp.LogErrorE("Hashid", err)
				return
			}
		}

		for k, v := range results1 {
			yapp.LogInfof("Hashid", "%d) %d = %s \n", k, *yapp.Hashid, v)
			if yapp.Savetocsv {

				res = []string{}
				//res = append(res,	fmt.Sprintf("%d", k))
				res = append(res, fmt.Sprintf("%d", *yapp.Hashid))
				res = append(res, v)

				err := writer.Write(res)
				if err != nil {
					yapp.LogError("Hashid", "Cannot save to files.csv")
					yapp.LogErrorE("Hashid", err)
					break
				}
			}

		}

	}
}
