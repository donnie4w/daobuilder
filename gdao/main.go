package main

import (
	"github.com/donnie4w/daobuilder/util"
	"github.com/donnie4w/simplelog/logging"
)

var gdaoversion = "1.1.0"

func init() {
	logging.SetOption(&logging.Option{Console: true, Format: logging.FORMAT_LEVELFLAG | logging.FORMAT_DATE | logging.FORMAT_TIME})
}
func main() {
	util.InitDB("gdao")
	fileBuilder()
}
