package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"encoding/csv"


	sl "github.com/eshu0/simplelogger"
	sli "github.com/eshu0/simplelogger/interfaces"
)

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func WalkDir(fds *DataStorage, Logger sli.ISimpleLogger, persist bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			Logger.LogErrorE("Visit", err)
			return nil
		}
		fexts := filepath.Ext(path)

		if strings.ToLower(fexts) != ".yaft" {

			//fd.FilePath = path
			abs, err := filepath.Abs(path)
			if err != nil {
				Logger.LogErrorE("Visit - Abs", err)
				return nil
			}

			fwn := FilenameWithoutExtension(abs)
			filename := filepath.Base(abs)
			fmt.Println("Filename:", filename)

			if filename[0] != '.' {
				fwn += ".yaft"
				hd := &HashData{}
				ok, _ := hd.CheckFileExists(fwn)
				if ok {
					data, ok := hd.LoadHashData(fwn, Logger)
					if ok {
						hd = data
					}
				} else {
					if !info.IsDir() {
						hr := HashRelationship{}

						if hr.GenHashData(Logger, abs, info.IsDir()) {
							fmt.Printf("%s %s\n", abs, hr.Hash.Data)
							hr.Path = abs
							fds.AddHashRelationship(&hr)
							if persist {
								hr.Hash.Save(fwn, Logger)
							}
						}
					}
				}
			} else {
				fmt.Printf("Hidden file %s \n", fwn)
			}

		}

		return nil
		/*
			} else {
				fmt.Printf("%s is directory\n", path)
				Logger.LogInfo("Visit", fmt.Sprintf("%s is directory", path))
				return nil
			}
		*/
	}
}

func main() {

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

	flag.Parse()

	savetocsav := false

	if savecsv != nil && *savecsv != "" {
		savetocsav = true
	}


	slog := sl.NewSimpleLogger(*filename, *session)

	// lets open a flie log using the session
	slog.OpenAllChannels()

	fds := &DataStorage{}
	fds.Filename = *dbname
	fds.Create(&slog)

	if inputdir != nil && *inputdir != "" {

		persist := (cache != nil && *cache != "")

		err := filepath.Walk(*inputdir, WalkDir(fds, &slog, persist))
		if err != nil {
			panic(err)
		}
	}

	if list != nil && *list != "" {
		fmt.Println("Listing all ")
		results := fds.GetAllHashData()

		fmt.Println("Hashs: ")
		for _, hd := range results {
			fmt.Println(hd)
		}

		fmt.Println("Relations: ")
		results2 := fds.GetAllHashRelationships()
		for _, hr := range results2 {
			fmt.Println(hr)
		}

	}

	if clear != nil && *clear != "" {
		fds.Clear()

	}

	if dupes != nil && *dupes != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if(savetocsav){
			file, err = os.Create("results.csv")
			if err != nil {
				slog.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
				return
			}
			defer file.Close()
		}

		fmt.Println("Duplicates: ")
		limitcount := -1
		if limit != nil && *limit > 0 {
			limitcount = *limit
		}

		results1 := fds.GetDuplicateHashes(limitcount)

		SaveDuplicates("./results.json",&slog,results1)
		if(savetocsav){
				writer = csv.NewWriter(file)
				defer writer.Flush()
		}

		for k, v := range results1 {
			fmt.Println("Key ", k)
			for _, hr := range v {
				fmt.Println(hr)

				if savetocsav {
					err := writer.Write(hr.CSV())
					if err != nil {
						slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
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

			if(savetocsav){
				file, err = os.Create("ids.csv")
				if err != nil {
					slog.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
					return
				}

				defer file.Close()
			}

			fmt.Println("Duplicates: ")
			limitcount := -1
			if limit != nil && *limit > 0 {
				limitcount = *limit
			}

			results1 := fds.GetDuplicateHashIds(limitcount)

			if(savetocsav){
					writer = csv.NewWriter(file)
					defer writer.Flush()
			}

			var res []string
			if savetocsav {

				res = []string{}
				//
				res = append(res,	"Hash Id")
				res = append(res,	"Count")

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
					res = append(res,	fmt.Sprintf("%d", v.HashId))
					res = append(res,	fmt.Sprintf("%d", v.Count))

					err := writer.Write(res)
					if err != nil {
						slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
						break
					}
				}


			}


	}
}
