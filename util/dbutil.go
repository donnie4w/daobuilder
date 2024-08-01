// Copyright (c) 2024, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/gdaobuilder

package util

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/donnie4w/gdao"
	"github.com/donnie4w/gdao/base"
	"log"

	//_ "github.com/alexbrainman/odbc"
	_ "github.com/denisenkom/go-mssqldb"
	goutil "github.com/donnie4w/gofer/util"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	//_ "github.com/ibmdb/go_ibm_db"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/nakagami/firebirdsql"
	_ "github.com/thda/tds"
	_ "github.com/vertica/vertica-sql-go"
	"os"
	"strings"
)

func openDB(driver string, config *ConfBean) (err error) {
	dataSourceName := ""
	switch strings.ToLower(driver) {
	case "mysql", "mariadb":
		dataSourceName = config.DbUser + ":" + config.DbPwd + "@tcp(" + config.DbHost + ":" + fmt.Sprint(config.DbPort) + ")/" + config.DbName
		DB, err = sql.Open("mysql", dataSourceName)
	case "postgresql":
		dataSourceName = "host=" + config.DbHost + " port=" + fmt.Sprint(config.DbPort) + " user=" + config.DbUser + " password=" + config.DbPwd + " dbname=" + config.DbName + " sslmode=disable"
		DB, err = sql.Open("postgres", dataSourceName)
	case "sqlite":
		dataSourceName = config.DbName
		DB, err = sql.Open("sqlite3", dataSourceName)
	case "sqlserver":
		dataSourceName = fmt.Sprint(`user="`, config.DbUser, `" password="`, config.DbPwd, `" connectString="`, config.DbHost, ":", config.DbPort, "/", config.DbName, `"`)
		DB, err = sql.Open("sqlserver", dataSourceName)
	case "oracle":
		dataSourceName = fmt.Sprint("server=", config.DbHost, ";port=", config.DbPort, ";userid=", config.DbUser, ";password=", config.DbPwd, ";database=", config.DbName)
		DB, err = sql.Open("godror", dataSourceName)
	case "h2":
		err = errors.New("H2 database files should be created using Jdao")
	case "db2":
		//"HOSTNAME=localhost;PORT=50000;DATABASE=mydb;UID=myuser;PWD=mypassword;"
		dataSourceName = fmt.Sprintf("HOSTNAME=%s;PORT=%d;DATABASE=%s;UID=%s;PWD=%s;", Config.DbHost, Config.DbPort, config.DbName, config.DbUser, config.DbPwd)
		DB, err = sql.Open("go_ibm_db", dataSourceName)
	case "sybase":
		dataSourceName = fmt.Sprintf("tds://%s:%s@%s:%d/%s", Config.DbUser, Config.DbPwd, config.DbHost, config.DbPort, config.DbName)
		DB, err = sql.Open("tds", dataSourceName)
	case "derby":
		//dataSourceName = fmt.Sprintf("jdbc:derby://%s:%d/%s;user=%s;password=%s", config.DbHost, config.DbPort, config.DbName, Config.DbUser, Config.DbPwd)
		//DB, err = sql.Open("oci8", dataSourceName)
		err = errors.New("derby database files should be created using Jdao")
	case "firebird":
		//connStr := "user:password@servername[:port_number]/database_name_or_file[?params1=value1[&param2=value2]...]"
		dataSourceName = fmt.Sprintf("%s:%s@%s[:%d]/%s", config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName)
		DB, err = sql.Open("firebirdsql", dataSourceName)
	case "ingres":
		dataSourceName = fmt.Sprintf("DSN=%s;UID=%s;PWD=%s;", config.DbName, config.DbUser, config.DbPwd)
		DB, err = sql.Open("odbc", dataSourceName)
	case "greenplum":
		dataSourceName = "host=" + config.DbHost + " port=" + fmt.Sprint(config.DbPort) + " user=" + config.DbUser + " password=" + config.DbPwd + " dbname=" + config.DbName + " sslmode=disable"
		DB, err = sql.Open("postgres", dataSourceName)
	case "teradata":
		dataSourceName = fmt.Sprintf("DSN=%s;UID=%s;PWD=%s;", config.DbName, config.DbUser, config.DbPwd)
		DB, err = sql.Open("odbc", dataSourceName)
	case "netezza":
		dataSourceName = fmt.Sprintf("DSN=%s;UID=%s;PWD=%s;", config.DbName, config.DbUser, config.DbPwd)
		DB, err = sql.Open("odbc", dataSourceName)
	case "vertica":
		dataSourceName = fmt.Sprintf("vertica://%s:%s@%s:%d/%s?connection_load_balance=1", config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName)
		DB, err = sql.Open("vertica", dataSourceName)
	default:
		err = fmt.Errorf("Unsupported driver: %s", driver)
	}
	log.Println("connect to ", driver, "[", dataSourceName, "]")
	if err == nil {
		if err = DB.Ping(); err != nil {
			log.Println("Ping DB failed:", err)
		}
	}
	return
}

func InitDB(daotype string) {
	daojson := ""
	log.Println("builder dao for ", daotype)
	flag.StringVar(&daojson, "c", daotype+".json", "configuration file of "+daotype+" in json")
	flag.Parse()
	if bs, err := goutil.ReadFile(daojson); err == nil {
		if Config, err = goutil.JsonDecode[*ConfBean](bs); err == nil {
			if err := openDB(Config.DbType, Config); err != nil {
				log.Println("open DB error:", err)
				os.Exit(0)
			}
			gdao.Init(DB, dbType(Config.DbType))
		} else {
			log.Println("decode config error:", err)
			os.Exit(0)
		}
	} else {
		log.Println("read config error:", err)
		os.Exit(0)
	}
	return
}

func dbType(s string) base.DBType {
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
	case "db2":
		return gdao.DB2
	case "sybase":
		return gdao.SYBASE
	case "derby":
		return gdao.DERBY
	case "firebird":
		return gdao.FIREBIRD
	case "ingres":
		return gdao.INGRES
	case "greenplum":
		return gdao.GREENPLUM
	case "teradata":
		return gdao.TERADATA
	case "netezza":
		return gdao.NETEZZA
	case "vertica":
		return gdao.VERTICA
	default:
		return gdao.MYSQL
	}
}

func Showtables() (r []string) {
	r = make([]string, 0)
	if databeans, err := gdao.ExecuteQueryBeans(showtable(Config.DbType)); err == nil {
		for _, bean := range databeans {
			r = append(r, bean.FieldByIndex(0).ValueString())
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
func (s showtables) DB2() string       { return "SELECT TABNAME FROM SYSCAT.TABLES" }
func (s showtables) Sybase() string    { return "SELECT name FROM sysobjects WHERE type = 'U'" }
func (s showtables) Derby() string {
	return "SELECT TABLENAME FROM SYS.SYSTABLES WHERE TABLETYPE = 'T'"
}
func (s showtables) Firebird() string {
	return "SELECT RDB$RELATION_NAME FROM RDB$RELATIONS WHERE RDB$VIEW_BLR IS NULL AND RDB$SYSTEM_FLAG = 0"
}
func (s showtables) Ingres() string {
	return "SELECT table_name FROM iitables WHERE table_owner != 'ingres'"
}
func (s showtables) Greenplum() string {
	return "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
}
func (s showtables) Teradata() string {
	return fmt.Sprintf("SELECT TableName FROM DBC.TablesV WHERE DatabaseName = '%s'", Config.DbName)
}
func (s showtables) Netezza() string { return "SELECT tablename FROM _v_table" }
func (s showtables) Vertica() string {
	return "SELECT table_name FROM v_catalog.tables WHERE table_schema='public'"
}
func (s showtables) H2() string {
	return "SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='PUBLIC'"
}

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
	case strings.EqualFold(s, "DB2"):
		return showtables(0).DB2()
	case strings.EqualFold(s, "SYBASE"):
		return showtables(0).Sybase()
	case strings.EqualFold(s, "DERBY"):
		return showtables(0).Derby()
	case strings.EqualFold(s, "FIREBIRD"):
		return showtables(0).Firebird()
	case strings.EqualFold(s, "INGRES"):
		return showtables(0).Ingres()
	case strings.EqualFold(s, "GREENPLUM"):
		return showtables(0).Greenplum()
	case strings.EqualFold(s, "TERADATA"):
		return showtables(0).Teradata()
	case strings.EqualFold(s, "NETEZZA"):
		return showtables(0).Netezza()
	case strings.EqualFold(s, "VERTICA"):
		return showtables(0).Vertica()
	default:
		return showtables(0).Mysql()
	}
}
