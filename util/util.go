// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package util

import (
	"log"
	"os"
	"strings"
)

func ToUpperFirstLetter(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
}

func CreeteInitfile(filename string) {
	if isFileExist(filename) {
		log.Println(filename, " is exist")
		return
	}
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write([]byte(jsonData))
	if err != nil {
		panic(err)
	}
	log.Println(filename, "  is created successfully")
}

var jsonData = `{
  "dbtype": "mysql",
  "dbhost": "localhost",
  "dbport": 3306,
  "dbname": "hstest",
  "dbuser": "root",
  "dbpwd": "123456",
  "package": "dao",
  "table": [""],
  "table_alias": [{"table": "","alias": ""}],
  "table_except": [""]
}`

func isFileExist(path string) (_r bool) {
	if path != "" {
		_, err := os.Stat(path)
		_r = err == nil || os.IsExist(err)
	}
	return
}
