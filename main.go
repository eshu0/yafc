package main

import (
	yaft "github.com/eshu0/yaft/pkg"
)

func main() {

	yapp := yaft.YAFTApp{}

	// parse the input flags
	yapp.LogInfo("Parsing Flags")
	yapp.ParseFlags()

	// create the database
	yapp.LogInfo("Creating Database")
	yapp.Create()

	yapp.LogInfo("Processing")
	yapp.Process()

	yapp.FinishLogging()
}
