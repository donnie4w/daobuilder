package main

import (
	"fmt"
	"github.com/donnie4w/daobuilder/util"
	. "github.com/donnie4w/gdao/gdaoBuilder"
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
		packageName = strings.ReplaceAll(packageName, "/", ".")
		packageName = strings.ReplaceAll(packageName, "\\", ".")
		packagePath := strings.ReplaceAll(packageName, ".", "/")

		func() {
			ok := false
			if str := createFile(tableName, tableAlias, packageName); str != "" {
				fileName := packagePath + "/" + util.ToUpperFirstLetter(tableAlias) + ".java"
				os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
				if f, err := os.Create(fileName); err == nil {
					defer f.Close()
					f.WriteString(str)
					log.Println("[successfully created dao file]", "[table:", tableName, "]["+fileName+"]")
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
	if tableBean, err := GetTableBean(tableName, util.DB); err == nil {
		dbtype, dbname := util.Config.DbType, util.Config.DbName
		return createDaoFile(dbtype, dbname, tableName, tableAlias, packageName, tableBean)
	} else {
		log.Println("[createFile] GetTableBean failed", err)
	}
	return ""
}

func createDaoFile(dbtype, dbname, tableName, tableAlias string, packageName string, tableBean *TableBean) (r string) {
	datetime := time.Now().Format(time.DateTime)
	ua := util.ToUpperFirstLetter
	structName := ua(tableAlias)

	timePackage := func() string {
		for _, bean := range tableBean.Fieldlist {
			if goType(bean.FieldType, bean.FieldTypeName) == "Date" {
				return "import java.util.Date;"
			}
		}
		return ""
	}()

	bigDecimalPackage := func() string {
		for _, bean := range tableBean.Fieldlist {
			if goType(bean.FieldType, bean.FieldTypeName) == "BigDecimal" {
				return "import java.math.BigDecimal;"
			}
		}
		return ""
	}()

	objectPackage := func() string {
		for _, bean := range tableBean.Fieldlist {
			rtype := goType(bean.FieldType, bean.FieldTypeName)
			if rtype == "BigDecimal" || rtype == "String" || rtype == "Date" {
				return `import java.util.Objects;`
			}
		}
		return ""
	}

	arrayPackage := func() string {
		for _, bean := range tableBean.Fieldlist {
			rtype := goType(bean.FieldType, bean.FieldTypeName)
			if rtype == "byte[]" {
				return `import java.util.Arrays;`
			}
		}
		return ""
	}

	r = `/*
 * Copyright (c) 2024, donnie4w <donnie4w@gmail.com> All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * github.com/donnie4w/jdao
 */

package ` + packageName + `;

import io.github.donnie4w.jdao.base.Fields;
import io.github.donnie4w.jdao.base.Table;
import io.github.donnie4w.jdao.base.Util;
import io.github.donnie4w.jdao.handle.JdaoException;
import io.github.donnie4w.jdao.util.Serializer;
import java.util.Map;
import java.util.HashMap;
` + arrayPackage() + `
` + objectPackage() + `
` + timePackage + `
` + bigDecimalPackage + `
/**
 * dbtype:` + dbtype + ` ,database:` + dbname + ` ,table:` + tableName + `
 *
 * @version jdao version ` + jdaoverion + `
 * @date ` + datetime + ` 
 */
public class ` + structName + ` extends Table<` + structName + `> {

	private static final long serialVersionUID = 6118074828633154000L;

	private final static String TABLENAME_ = "` + tableName + `";
`
	fieldstr := ""
	fieldToStr := ""
	for i, bean := range tableBean.Fieldlist {
		log.Println(bean)
		fieldName := encodeFieldname(bean.FieldName)
		fieldstr += strings.ToUpper(bean.FieldName)
		fieldToStr += `"` + bean.FieldName + `:" + ` + fieldName
		if i < len(tableBean.Fieldlist)-1 {
			fieldstr += ","
			fieldToStr += ` + " , " + `
		}
		r += `
	public final static Fields<` + structName + `> ` + strings.ToUpper(bean.FieldName) + ` = new Fields("` + bean.FieldName + `");`
	}

	r += `
	
	public ` + ua(structName) + `() {
		super(TABLENAME_, ` + ua(structName) + `.class);
		super.initFields(` + fieldstr + `);
	}

	public ` + ua(structName) + `(String tableName) {
		super(tableName, ` + ua(structName) + `.class);
		super.initFields(` + fieldstr + `);
	}

	@Override
	public void toJdao() {
		super.init(` + structName + `.class);
	}
`

	for _, bean := range tableBean.Fieldlist {
		rtype := goType(bean.FieldType, bean.FieldTypeName)
		fieldname := encodeFieldname(bean.FieldName)
		r += `
	private ` + rtype + ` ` + fieldname + `;`
	}
	r += "\n"
	for _, bean := range tableBean.Fieldlist {
		rtype := goType(bean.FieldType, bean.FieldTypeName)
		fieldname := encodeFieldname(bean.FieldName)
		r += `
	public ` + rtype + ` get` + ua(bean.FieldName) + `() {
		return this.` + fieldname + `;
	}

	public void set` + ua(bean.FieldName) + `(` + rtype + ` ` + fieldname + `) {
		fieldPut(` + strings.ToUpper(bean.FieldName) + `, ` + fieldname + `);
		this.` + fieldname + ` = ` + fieldname + `;
	}
`
	}

	r += `
	@Override
	public String toString() {
		return ` + fieldToStr + `;
	}
`

	copy := `
	@Override
	public ` + structName + ` copy(` + structName + ` h) {`
	r = r + copy
	for _, bean := range tableBean.Fieldlist {
		r = r + `
		this.set` + ua(bean.FieldName) + `(h.get` + ua(bean.FieldName) + `());`
	}
	r = r + `
		return this;
	}
`

	r += `
	@Override
	public void scan(String fieldname, Object obj) throws JdaoException {
		try {
			switch (fieldname) {`
	for _, bean := range tableBean.Fieldlist {
		rtype := goType(bean.FieldType, bean.FieldTypeName)
		if bean.FieldType.Kind() == reflect.Slice && bean.FieldType.Elem().Kind() == reflect.Uint8 {
			rtype = "bytes"
		}
		r += `
				case "` + bean.FieldName + `":
					set` + ua(bean.FieldName) + `(Util.as` + ua(rtype) + `(obj));
 					break;`

	}

	r += `
			}
		} catch (Exception e) {
			throw new JdaoException(e);
		}`

	r += `
	}
`

	equalstr := ""
	hashCodeStr := ""
	structName2 := "_" + structName
	for i, bean := range tableBean.Fieldlist {
		fieldName := encodeFieldname(bean.FieldName)
		rtype := goType(bean.FieldType, bean.FieldTypeName)
		if rtype == "BigDecimal" || rtype == "String" || rtype == "Date" {
			equalstr += `Objects.equals(` + fieldName + `, ` + structName2 + `.` + fieldName + `)`
			hashCodeStr += fieldName
		} else if rtype == "byte[]" {
			equalstr += `Objects.deepEquals(` + fieldName + `, ` + structName2 + `.` + fieldName + `)`
			hashCodeStr += `Arrays.hashCode(` + fieldName + `)`
		} else {
			equalstr += fieldName + " == " + structName2 + "." + fieldName
			hashCodeStr += fieldName
		}

		if i < len(tableBean.Fieldlist)-1 {
			equalstr += " && "
			hashCodeStr += " , "
		}

	}

	seriastr := `
	@Override
	public byte[] encode() {
		Map<String, Object> map = new HashMap();`
	r = r + seriastr
	for _, bean := range tableBean.Fieldlist {
		r = r + `
		map.put("` + bean.FieldName + `", this.get` + ua(bean.FieldName) + `());`
	}

	r = r + `
		return Serializer.encode(map);
	}

	@Override
	public ` + structName + ` decode(byte[] bs) throws JdaoException {
		Map<String, Object> map = Serializer.decode(bs);
		if (map != null) {
			for (Map.Entry<String, Object> entry : map.entrySet()) {
				scan(entry.getKey(), entry.getValue());
			}
		}
		return this;
	}
`

	r += `
	@Override
	public boolean equals(Object o) {
		if (this == o) return true;
		if (o == null || getClass() != o.getClass()) return false;
        ` + structName + ` ` + structName2 + ` = (` + structName + `) o;
		return ` + equalstr + `;
	}

	@Override
	public int hashCode() {
		return Objects.hash(` + hashCodeStr + `);
	}
}`
	return
}
