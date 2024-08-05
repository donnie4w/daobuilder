// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package main

import (
	"github.com/donnie4w/daobuilder/util"
	"github.com/donnie4w/gdao/gdaoBuilder"
	"log"
)

func fileBuilder() error {
	tables := util.Config.Table

	if len(tables) == 0 || tables[0] == "" {
		tables = util.Showtables()
	}
	ts := make([]string, 0)
	for _, tableName := range tables {
		if util.Config.IsExcept(tableName) {
			continue
		}
		ts = append(ts, tableName)
	}
	log.Println("[tables]", ts)
	for _, tableName := range ts {
		tableAlias := util.Config.GetAlias(tableName)
		gdaoBuilder.BuildWithAlias(tableName, tableAlias, util.Config.DbType, util.Config.DbName, util.Config.Package, util.DB)
	}
	return nil
}
