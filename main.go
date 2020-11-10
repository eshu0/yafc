package main

import (
	yaft "github.com/eshu0/yaft/pkg"
)

func main() {

	yapp := yaft.YAFTApp{}

	// parse the input flags
	yapp.ParseFlags()

	// create the database
	yapp.Create()

}
