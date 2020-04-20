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

func WalkDir(fds *DataStorage, Logger sli.ISimpleLogger) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			Logger.LogErrorE("Visit", err)
			return nil
		}
		fexts := filepath.Ext(path)

		if strings.ToLower(fexts) != ".yaft" {

			//fd.FilePath = path
			abs, err := filepath.Abs(path)
			if err == nil {
				fmt.Println("Absolute:", abs)
			}

			fwn := FilenameWithoutExtension(abs)
			if fwn[0] != '.' {
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
							hr.Hash.Save(fwn, Logger)
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
	list := flag.String("list", "", "")

	flag.Parse()

	slog := sl.NewSimpleLogger(*filename, *session)

	// lets open a flie log using the session
	slog.OpenAllChannels()

	fds := &DataStorage{}
	fds.Filename = *dbname
	fds.Create()

	if inputdir != nil && *inputdir != "" {
		err := filepath.Walk(*inputdir, WalkDir(fds, &slog))
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
		/*
			fmt.Println("Files: ")
			results1 := fds.GetAllFileData()
			for _, fd := range results1 {
				fmt.Println(fd)
			}
		*/
		fmt.Println("Relations: ")
		results2 := fds.GetAllHashRelationships()
		for _, hr := range results2 {
			fmt.Println(hr)
		}

	}
}
