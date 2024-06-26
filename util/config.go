// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package util

const (
	_VER = "1.0.2"
)

type ConfBean struct {
	DbType      string        `json:"dbtype"`
	DbHost      string        `json:"dbhost"`
	DbPort      int           `json:"dbport"`
	DbName      string        `json:"dbname"`
	DbUser      string        `json:"dbuser"`
	DbPwd       string        `json:"dbpwd"`
	Package     string        `json:"package"`
	Table       []string      `json:"table"`
	TableExcept []string      `json:"table_except"`
	TableAlias  []Table_alias `json:"table_alias"`
}

func (c *ConfBean) GetAlias(tablename string) string {
	if len(c.TableAlias) > 0 {
		for _, v := range c.TableAlias {
			if v.Table == tablename {
				return v.Alias
			}
		}
	}
	return tablename
}

func (c *ConfBean) IsExcept(tablename string) bool {
	if len(c.TableExcept) > 0 {
		for _, v := range c.TableExcept {
			if v == tablename {
				return true
			}
		}
	}
	return false
}

type Table_alias struct {
	Table string `json:"table"`
	Alias string `json:"alias"`
}
