package yaft

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"bufio"

	sli "github.com/eshu0/simplelogger/interfaces"
)

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}


func CompareDirectory(fds *DataStorage, Logger sli.ISimpleLogger, die *bool, reader *bufio.Reader) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {

		if err != nil {
			Logger.LogErrorE("Visit", err)
			return nil
		}

		//fd.FilePath = path
		abs, err := filepath.Abs(path)
		if err != nil {
			Logger.LogErrorE("Visit - Abs", err)
			return nil
		}

		if !info.IsDir() {
			hr := HashRelationship{}

			if hr.GenHashData(Logger, abs, info.IsDir()) {
				fmt.Printf("%s %s\n", abs, hr.Hash.Data)
				hr.Path = abs
				res := fds.FindHashData(hr.Hash.Data)
				for _, v := range res {
					hrs := fds.GetHashRelationshipByHash(v.ID)
					for _, hr := range hrs {
						fmt.Println(hr.Path)
					}
					if die != nil && *die && len(hrs) > 0{
						fmt.Printf("Delete file %s: \n",path)
						text, _ := reader.ReadString('\n')
						if strings.Compare("yes", text) == 0 {
							err := os.Remove(path)
							if err != nil {
								fmt.Println(err)
								return nil
							}else{
								fmt.Printf("deleted file %s",path)
							}
						}else{
							fmt.Printf("Not deleting file %s",path)
						}
					}
				}
			}
		}

		return nil
	}
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
