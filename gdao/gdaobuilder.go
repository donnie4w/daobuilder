// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package main

import (
	"database/sql"
	"fmt"
	"github.com/donnie4w/daobuilder/util"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
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
		packageName := util.Config.Package
		func() {
			ok := false
			if str := createFile(tableName, tableAlias, packageName); str != "" {
				fileName := packageName + "/" + tableAlias + ".go"
				os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
				if f, err := os.Create(fileName); err == nil {
					defer f.Close()
					f.WriteString(str)
					log.Println("[successfully created dao file]", createlog(tableName, tableAlias))
					ok = true
				}
			}
			if !ok {
				log.Println("[failed to created dao file]", createlog(tableName, tableAlias))
			}
		}()
	}
	return nil
}

func createlog(tableName, tableAlias string) string {
	if tableAlias != "" && tableAlias != tableName {
		return fmt.Sprint("["+tableName+" ]As[", tableAlias, " ]")
	}
	return tableName
}

func createFile(tableName, tableAlias string, packageName string) string {
	tableBean := util.TableInfo(tableName)
	dbtype, dbname := util.Config.DbType, util.Config.DbName
	return createDaoFile(dbtype, dbname, tableName, tableAlias, packageName, tableBean)
}

func createDaoFile(dbtype, dbname, tableName, tableAlias string, packageName string, tableBean *util.TableBean) string {
	datetime := time.Now().Format(time.DateTime)
	ua := util.ToUpperFirstLetter
	structName := ua(tableAlias)

	r := `// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// github.com/donnie4w/gdao
// datetime :` + datetime + `
// gdao version ` + gdaoversion + `
// ` + dbtype + `
// database:` + dbname + ` ,tablename:` + tableName + `

package ` + packageName + `

import (
	"fmt"
	"github.com/donnie4w/gdao"
	` + func() string {
		b := false
		for _, bean := range tableBean.Fieldlist {
			if bean.FieldType == reflect.TypeOf(sql.NullTime{}) {
				b = true
				break
			}
		}
		if b {
			return "\"time\""
		} else {
			return ""
		}
	}()
	r = r + `
)`

	for _, bean := range tableBean.Fieldlist {
		rtype := goPtrType(bean.FieldType)
		structField := strings.ToLower(tableAlias) + "_" + ua(bean.FieldName)
		s := `
type ` + structField + ` struct {
	gdao.Field
	fieldName  string
	FieldValue ` + rtype + `
}

func (t *` + structField + `) Name() string {
	return t.fieldName
}

func (t *` + structField + `) Value() any {
	return t.FieldValue
}
`
		r = r + s
	}

	r = r + `
type ` + structName + ` struct {
	gdao.Table
`
	for _, bean := range tableBean.Fieldlist {
		s := `
	` + ua(bean.FieldName) + `		*` + strings.ToLower(tableAlias) + "_" + ua(bean.FieldName)
		r = r + s
	}
	r = r + `
}
`
	mustptr := func(t reflect.Type, s string) string {
		if mustPtr(t) {
			return s
		} else {
			return ""
		}
	}

	for _, bean := range tableBean.Fieldlist {
		fieldName := ua(bean.FieldName)
		rtype := goType(bean.FieldType)
		s := `
func (u *` + structName + `) Get` + fieldName + `() (_r ` + rtype + `){
`
		s1 := `	if u.` + fieldName + `.FieldValue != nil {
		_r = ` + mustptr(bean.FieldType, "*") + `u.` + fieldName + `.FieldValue
	}`
		if mustPtr(bean.FieldType) {
			s = s + s1
		} else {
			s = s + `	_r = ` + mustptr(bean.FieldType, "*") + `u.` + fieldName + `.FieldValue`
		}
		s = s + `
	return
}

func (u *` + structName + `) Set` + fieldName + `(arg ` + rtype + `) *` + structName + `{
	u.Put0(u.` + fieldName + `.fieldName, arg)
	u.` + fieldName + `.FieldValue = ` + mustptr(bean.FieldType, "&") + `arg
	return u
}
`
		r = r + s
	}

	r = r + `

func (u *` + structName + `) Scan(fieldname string, value any) {
	switch fieldname {`

	for _, bean := range tableBean.Fieldlist {
		fieldName := ua(bean.FieldName)
		rtype := goType(bean.FieldType)
		var s string
		if bean.FieldType == reflect.TypeOf(sql.NullTime{}) {

			s = `
	case "` + bean.FieldName + `":
		if t, err := gdao.AsTime(value); err == nil {
		u.Set` + fieldName + `(t)
		}`
		} else if bean.FieldType.Kind() == reflect.Slice && bean.FieldType.Elem().Kind() == reflect.Uint8 {

			s = `
	case "` + bean.FieldName + `":
		u.Set` + fieldName + `(gdao.AsBytes(value))`
		} else {

			s = `
	case "` + bean.FieldName + `":
		u.Set` + fieldName + `(gdao.As` + ua(rtype) + `(value))`
		}
		r = r + s
	}
	r = r + `
	}
}
`
	columns := ""
	fields := ""
	fieldsString := ""
	for i, bean := range tableBean.Fieldlist {
		columns = columns + "t." + ua(bean.FieldName)
		fieldsString = fieldsString + "\"" + ua(bean.FieldName) + ":\"" + ",t.Get" + ua(bean.FieldName) + "()"
		fields = fields + ua(bean.FieldName) + ":" + CheckReserveKey(bean.FieldName)
		if i < len(tableBean.Fieldlist)-1 {
			columns = columns + ","
			fieldsString = fieldsString + `, ",",`
			fields = fields + ","
		}
	}

	selectfunc := `

func (t *` + structName + `) Selects(columns ...gdao.Column) (_r []*` + structName + `, err error) {
	if columns == nil {
		columns = []gdao.Column{` + columns + `}
	}
	databean, err := t.QueryBeans(columns...)
	if err != nil || len(databean) == 0 {
		return nil, err
	}
	_r = make([]*` + structName + `, 0)
	for _, beans := range databean {
		__` + structName + ` := New` + structName + `()
		for name, beans := range beans.FieldMapName {
			__` + structName + `.Scan(name, beans.Value())
		}
		_r = append(_r, __` + structName + `)
	}
	return
}

func (t *` + structName + `) Select(columns ...gdao.Column) (_r *` + structName + `, err error) {
	if columns == nil {
		columns = []gdao.Column{` + columns + `}
	}
	databean, err := t.QueryBean(columns...)
	if err != nil || databean == nil {
		return nil, err
	}
	_r = New` + structName + `()
	for name, beans := range databean.FieldMapName {
		_r.Scan(name, beans.Value())
	}
	return
}
`
	r = r + selectfunc

	r = r + `
func (t *` + structName + `) New() gdao.Scanner {
	return New` + structName + `()
}
`

	stringBody := `
func (t *` + structName + `) String() string {
	return fmt.Sprint(` + fieldsString + `)
}
`
	r = r + stringBody

	newfunc := `
func New` + structName + `(tablename ...string) (_r *` + structName + `) {
`
	for _, bean := range tableBean.Fieldlist {
		structField := strings.ToLower(tableAlias) + "_" + ua(bean.FieldName)
		varfield := CheckReserveKey(bean.FieldName)
		s := `
	` + varfield + ` := &` + structField + `{fieldName: "` + bean.FieldName + `"}
	` + varfield + `.Field.FieldName = "` + bean.FieldName + `"
`
		newfunc = newfunc + s
	}

	newfunc = newfunc + `
	_r = &` + structName + `{` + fields + `}
	s := "` + tableName + `"
	if len(tablename) > 0 && tablename[0] != "" {
		s = tablename[0]
	}
	_r.Init(s)
	return
}
`
	r = r + newfunc
	return r
}
