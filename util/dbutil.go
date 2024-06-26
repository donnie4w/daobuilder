// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package util

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/donnie4w/gdao"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func openDB(driver string, config *ConfBean) (err error) {
	switch strings.ToLower(driver) {
	case "mysql", "mariadb":
		dataSourceName := config.DbUser + ":" + config.DbPwd + "@tcp(" + config.DbHost + ":" + fmt.Sprint(config.DbPort) + ")/" + config.DbName
		logging.Info("connect to ", driver, "[", dataSourceName, "]")
		if DB, err = sql.Open("mysql", dataSourceName); err != nil {
			logging.Error("open DB error:", err)
		}
	case "postgresql":
		dataSourceName := "host=" + config.DbHost + " port=" + fmt.Sprint(config.DbPort) + " user=" + config.DbUser + " password=" + config.DbPwd + " dbname=" + config.DbName + " sslmode=disable"
		logging.Info("connect to ", driver, "[", dataSourceName, "]")
		if DB, err = sql.Open("postgres", dataSourceName); err != nil {
			logging.Error("open DB error:", err)
		}
	case "sqlite":
		dataSourceName := config.DbName
		if DB, err = sql.Open("sqlite3", dataSourceName); err != nil {
			logging.Error("open DB error:", err)
		}
	case "sqlserver":
		dataSourceName := fmt.Sprint(`user="`, config.DbUser, `" password="`, config.DbPwd, `" connectString="`, config.DbHost, ":", config.DbPort, "/", config.DbName, `"`)
		if DB, err = sql.Open("sqlserver", dataSourceName); err != nil {
			logging.Error("open DB error:", err)
		}
	case "oracle":
		dataSourceName := fmt.Sprint("server=", config.DbHost, ";port=", config.DbPort, ";userid=", config.DbUser, ";password=", config.DbPwd, ";database=", config.DbName)
		if DB, err = sql.Open("godror", dataSourceName); err != nil {
			logging.Error("open DB error:", err)
		}
	default:
		err = fmt.Errorf("Unsupported driver: %s", driver)
	}
	return
}

func InitDB(daotype string) {
	daojson := ""
	logging.Info("builder dao for ", daotype)
	flag.StringVar(&daojson, "c", daotype+".json", "configuration file of "+daotype+" in json")
	flag.Parse()
	if bs, err := goutil.ReadFile(daojson); err == nil {
		if Config, err = goutil.JsonDecode[*ConfBean](bs); err == nil {
			if err := openDB(Config.DbType, Config); err != nil {
				logging.Error("open DB error:", err)
				os.Exit(0)
			}
			Gdao = gdao.NewGdao(DB, dbType(Config.DbType))
		} else {
			logging.Error("decode config error:", err)
			os.Exit(0)
		}
	} else {
		logging.Error("read config error:", err)
		os.Exit(0)
	}
	return
}

func dbType(s string) gdao.DBType {
	switch strings.ToLower(s) {
	case "mysql":
		return gdao.MYSQL
	case "postgresql":
		return gdao.POSTGRESQL
	case "sqlite":
		return gdao.SQLITE
	case "oracle":
		return gdao.ORACLE
	case "sqlserver":
		return gdao.SQLSERVER
	case "mariadb":
		return gdao.MARIADB
	default:
		return gdao.MYSQL
	}
}

func Showtables() (r []string) {
	r = make([]string, 0)
	if databeans, err := Gdao.QueryBeans(showtable(Config.DbType)); err == nil {
		for _, bean := range databeans {
			r = append(r, bean.FieldMapIndex[0].ValueString())
		}
	}
	return
}

type showtables byte

func (s showtables) Mysql() string { return "SHOW TABLES" }
func (s showtables) PostgreSQL() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
}
func (s showtables) Mariadb() string   { return "SHOW TABLES" }
func (s showtables) Sqlite() string    { return "SELECT name FROM sqlite_master WHERE type = 'table' " }
func (s showtables) Oracle() string    { return "SELECT table_name FROM user_tables" }
func (s showtables) Sqlserver() string { return "SELECT name FROM sys.tables" }

func showtable(s string) string {
	switch {
	case strings.EqualFold(s, "mysql"):
		return showtables(0).Mysql()
	case strings.EqualFold(s, "Mariadb"):
		return showtables(0).Mariadb()
	case strings.EqualFold(s, "Sqlite"):
		return showtables(0).Sqlite()
	case strings.EqualFold(s, "Oracle"):
		return showtables(0).Oracle()
	case strings.EqualFold(s, "Sqlserver"):
		return showtables(0).Sqlserver()
	case strings.EqualFold(s, "PostgreSQL"):
		return showtables(0).PostgreSQL()
	default:
		return showtables(0).Mysql()
	}
}
