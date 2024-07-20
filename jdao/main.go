package main

import (
	"github.com/donnie4w/daobuilder/util"
	"os"
)

var jdaoverion = "2.0.1"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "init" {
		util.CreeteInitfile("jdao.json")
		return
	}
	util.InitDB("jdao")
	fileBuilder()
}
