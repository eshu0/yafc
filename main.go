package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	yaft "github.com/eshu0/yaft/pkg"
)

func main() {

	yapp := YAFTApp{}

	// parse the input flags
	yapp.ParseFlags()

	// create the database
	yapp.Create()

	if filetofind != nil && *filetofind != "" {
		reader := bufio.NewReader(os.Stdin)

		if deleteifexists != nil && *deleteifexists {
			fmt.Println("Delete if exists")
		}

		if yestoall != nil && *yestoall {
			fmt.Println("Yes to all")
		}

		err := filepath.Walk(*filetofind, yaft.CompareDirectory(fds, slog, deleteifexists, yestoall, reader))
		if err != nil {
			panic(err)
		}
	}

	if inputdir != nil && *inputdir != "" {

		persist := (cache != nil && *cache != "")

		err := filepath.Walk(*inputdir, yaft.WalkDir(fds, slog, persist))
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

	if clear != nil && *clear != "" {
		fds.Clear()
	}

	if dupes != nil && *dupes != "" {

		var file *os.File
		var writer *csv.Writer
		var err error

		if savetocsav {
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

		yaft.SaveDuplicates("./results.json", slog, results1)
		if savetocsav {
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
					slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
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
				slog.LogError("CreateCSV", fmt.Sprintf("Cannot create file%s", err.Error()))
				return
			}

			defer file.Close()
		}

		fmt.Println("Files: ")

		results1 := fds.GetFilesByHashId(int64(*hashid))

		if savetocsav {
			writer = csv.NewWriter(file)
			defer writer.Flush()
		}

		var res []string
		if savetocsav {

			res = []string{}
			//
			res = append(res, "Hash Id")
			res = append(res, "Files")

			err := writer.Write(res)
			if err != nil {
				slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
				return
			}
		}

		for k, v := range results1 {
			fmt.Printf("%d) %d = %s \n", k, *hashid, v)
			if savetocsav {

				res = []string{}
				//res = append(res,	fmt.Sprintf("%d", k))
				res = append(res, fmt.Sprintf("%d", *hashid))
				res = append(res, v)

				err := writer.Write(res)
				if err != nil {
					slog.LogError("CreateCSV", fmt.Sprintf("Cannot write to file %s", err.Error()))
					break
				}
			}

		}

	}

}
