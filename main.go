package main

import (
	yaft "github.com/eshu0/yaft/pkg"
)

func main() {

	yapp := yaft.YAFTApp{}

	yapp.LogInfo("Parsing Flags")

	// parse the input flags
	yapp.ParseFlags()

	yapp.LogInfo("Creating Database")
	// create the database
	yapp.Create()

	yapp.LogInfo("Processing")
	yapp.Process()
}
