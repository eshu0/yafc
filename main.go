package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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
	dbname := flag.String("dbname", "./yaft.db", "Database defaults to ./yaft.db")
	inputdir := flag.String("path", "", "")
	cache := flag.String("cache", "", "")
	list := flag.String("list", "", "")
	clear := flag.String("clear", "", "")
	dupes := flag.String("dupes", "", "")

	flag.Parse()

	slog := sl.NewSimpleLogger(*filename, *session)

	// lets open a flie log using the session
	slog.OpenAllChannels()

	fds := &DataStorage{}
	fds.Filename = *dbname
	fds.Create()

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
		fmt.Println("Duplicates: ")
		results1 := fds.GetDuplicateHashes()
		for k, v := range results1 {
			fmt.Println("Key ", k)
			for _, hr := range v {
				fmt.Println(hr)
			}
		}
	}

}
