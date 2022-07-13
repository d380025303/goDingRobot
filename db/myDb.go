package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

func InitSqlite(location string) {
	log.Println("init sqlite3...")
	myDb, err := sql.Open("sqlite3", location)
	if err != nil {
		panic(err)
	}
	err = myDb.Ping()
	if err != nil {
		panic("连接数据库失败: ")
	}
	db = myDb
	InsertOrUpdateSqlByStmt(`CREATE TABLE IF NOT EXISTS message ( 
		   id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		   content TEXT not NULL,
		   third_message_id TEXT not null,
           from_user_id TEXT not null);`)
}

func GetDb() *sql.DB {
	myDb := db
	if myDb == nil {
		panic("db not init...")
	}
	return myDb
}

func InsertOrUpdateSqlByStmt(sqlStr string, args ...interface{}) {
	myDb := GetDb()
	stmt, _ := myDb.Prepare(sqlStr)
	defer stmt.Close()
	stmt.Exec(args...)
}
