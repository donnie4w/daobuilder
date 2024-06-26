// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdao

package util

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/gdao"
	"reflect"
	"time"
)

var DB *sql.DB
var Gdao *gdao.Gdao
var Config *ConfBean

type TableBean struct {
	TableName string
	Fieldlist []*fieldBean
	Fieldmap  map[string]*fieldBean
}

func (t *TableBean) ContainTime() bool {
	for _, field := range t.Fieldlist {
		if field.FieldType == reflect.TypeOf(time.Time{}) {
			return true
		}
	}
	return false
}

func newTableBean() *TableBean {
	return &TableBean{Fieldlist: make([]*fieldBean, 0), Fieldmap: make(map[string]*fieldBean)}
}

type fieldBean struct {
	FieldName     string
	FieldIndex    int
	FieldType     reflect.Type
	FieldTypeName string
}

func (f *fieldBean) String() string {
	return fmt.Sprint(f.FieldName, ",", f.FieldIndex, ",", f.FieldType, ",", f.FieldTypeName)
}

func executeForTableInfo(sql string) (tb *TableBean, err error) {
	rows, err := DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tb = newTableBean()
	types, _ := rows.ColumnTypes()
	for i, columntype := range types {
		fb := &fieldBean{}
		fb.FieldTypeName = columntype.DatabaseTypeName()
		fb.FieldType = columntype.ScanType()
		fb.FieldName = columntype.Name()
		fb.FieldIndex = i
		tb.Fieldmap[columntype.Name()] = fb
		tb.Fieldlist = append(tb.Fieldlist, fb)
	}
	return
}

func TableInfo(tablename string) (tb *TableBean) {
	sql := fmt.Sprint("select * from ", tablename, " where 0=1")
	tb, _ = executeForTableInfo(sql)
	if tb != nil {
		tb.TableName = tablename
	}
	return
}
