package yaft

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

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
