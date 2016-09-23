package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"time"
	//"os"
)

var DB *sql.DB = nil

func typeOf(in interface{}) {
	fmt.Println(reflect.TypeOf(in))
}

func initDB() error {
	db, err := sql.Open("sqlite3", dbaddr)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(reflect.TypeOf(err))
	DB = db
	ver := verifyDB()
	switch ver {
	case -1:
		createDB()
	case dbversion:
		log.Println("Database version: ", dbversion)
	default:
		log.Fatal("Database version mismatch, upgrade required")
	}
	return err
}

func closeDB() {
	defer DB.Close()
}

func verifyDB() int {
	checkVersion := "select version from dbversion;"
	rows, err := DB.Query(checkVersion)
	if err != nil {
		return -1
	}
	rows.Next()
	var version int
	rows.Scan(&version)
	return version
}

func createDB() error {
	createTableLog := `
		create table log (
				timestamp integer not null,
				caller string not null,
				chatstring string not null,
				chatout string not null,
				systemcommand string,
				systemout string,
				systemerror string );`
	createTableHostaliases := `
		create table hostaliases (
				name string not null primary key,
				owner string not null,
				host string not null
		);`
	createTableDbversion := `
		create table dbversion (
				version integer
		);`
	insertDbversion := fmt.Sprintf(`
		insert into dbversion(version) values (%d);
		`, dbversion)
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createTableLog)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createTableHostaliases)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(createTableDbversion)
	if err != nil {
		log.Fatal(err)
	}
	_, err = DB.Exec(insertDbversion)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
	return err
}

func logToDB(logEntry dbLogEntry) {
	tx, err := DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	statement, err := DB.Prepare(`
					insert into log(
					timestamp,
					caller,
					chatstring,
					chatout,
					systemcommand,
					systemout,
					systemerror) values (?, ?, ?, ?, ?, ?, ?);`)
	defer statement.Close()
	syserr := ""
	if logEntry.systemerror != nil {
		syserr = fmt.Sprintf("%s", logEntry.systemerror)
	}
	statement.Exec(
		int32(time.Now().Unix()),
		logEntry.caller,
		logEntry.chatstring,
		logEntry.chatout,
		logEntry.systemcommand,
		logEntry.systemout,
		syserr)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
	tx.Commit()
}

func lastFromDB(params ...int) []dbLogEntry {
	nentries := 10
	if len(params) > 0 {
		nentries = params[0]
	}
	statement := fmt.Sprintf("select * from log order by timestamp desc limit %d", nentries)
	rows, err := DB.Query(statement)
	defer rows.Close()
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
	var entries []dbLogEntry
	for rows.Next() {
		var timestamp int
		var caller string
		var chatstring string
		var chatout string
		var systemcommand string
		var systemout string
		var systemerror string
		rows.Scan(&timestamp, &caller, &chatstring, &chatout, &systemcommand, &systemout, &systemerror)
		e := dbLogEntry{
			timestamp:     timestamp,
			caller:        caller,
			chatstring:    chatstring,
			chatout:       chatout,
			systemcommand: systemcommand,
			systemout:     systemout,
			systemerror:   errors.New(systemerror),
		}
		entries = append(entries, e)
	}
	return entries
}

func addAlias(alias dbAliasEntry) {
	statement := fmt.Sprintf(
		"insert into hostaliases(name,owner,host) values ('%s','%s','%s')",
		alias.name,
		alias.owner,
		alias.host)
	_, err := DB.Exec(statement)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
}

func delAlias(alias dbAliasEntry) {
	statement, err := DB.Prepare("delete from hostaliases where name='?' and owner='?'")
	defer statement.Close()
	statement.Exec(alias.name, alias.owner)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
}

func showAlias(alias dbAliasEntry) []dbAliasEntry {
	var entries []dbAliasEntry
	var statement string
	if alias.name != "" {
		statement = fmt.Sprintf("select * from hostaliases where name='%s' and owner='%s'", alias.name, alias.owner)
	} else {
		statement = fmt.Sprintf("select * from hostaliases where owner='%s'", alias.owner)
	}
	rows, err := DB.Query(statement)
	defer rows.Close()
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
		return entries
	}
	for rows.Next() {
		var name string
		var owner string
		var host string
		rows.Scan(&name, &owner, &host)
		e := dbAliasEntry{
			name:  name,
			owner: owner,
			host:  host,
		}
		entries = append(entries, e)
	}
	return entries
}

func updateAlias(alias dbAliasEntry) {
	aliases := showAlias(alias)
	switch {
	case len(aliases) == 0:
		addAlias(alias)
	case len(aliases) == 1:
		statement := fmt.Sprintf(
			"update hostaliases set host='%s' where name='%s' and owner='%s'",
			alias.host, alias.name, alias.owner,
		)
		_, err := DB.Exec(statement)
		if err != nil {
			log.Printf("%q: %s\n", err, statement)
		}
	case len(aliases) > 1:
		//shouldn't happen
	}
}
