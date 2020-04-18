package main

import (
	"bufio"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	sl "github.com/eshu0/simplelogger"
	sli "github.com/eshu0/simplelogger/interfaces"
)

func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

func WalkDir(Logger sli.ISimpleLogger) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Logger.LogErrorE("Visit", err)
		}

		isdir, err := IsDirectory(path)

		if !isdir {
			f, err := os.Open(path)

			if err != nil {
				Logger.LogErrorE("Visit", err)
				return nil
			} else {
				input := bufio.NewReader(f)

				hash := sha256.New()
				if _, err := io.Copy(hash, input); err != nil {
					Logger.LogErrorE("Visit", err)
				}
				sum := hash.Sum(nil)

				fmt.Printf("%s %x\n", path, sum)
				Logger.LogInfo("Visit", fmt.Sprintf("%s %x", path, sum))

				return nil
			}
		} else {
			fmt.Printf("%s is directory\n", path)
			Logger.LogInfo("Visit", fmt.Sprintf("%s is directory", path))
			return nil
		}

	}
}

func main() {

	filename := flag.String("logfile", "yaft.log", "Filename out - defaults to yaft.log")
	session := flag.String("sessionid", "123", "Session - defaults to 123")
	inputdir := flag.String("path", "", "")

	flag.Parse()

	slog := sl.NewSimpleLogger(*filename, *session)

	// lets open a flie log using the session
	slog.OpenAllChannels()

	err := filepath.Walk(*inputdir, WalkDir(&slog))
	if err != nil {
		panic(err)
	}

}
