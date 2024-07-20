package main

import (
	"github.com/donnie4w/daobuilder/util"
	"log"
	"os"
)

var gdaoversion = "1.1.0"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "init" {
		util.CreeteInitfile("mysql.json")
		return
	}
	log.Println(os.Args)
	util.InitDB("gdao")
	fileBuilder()
}
