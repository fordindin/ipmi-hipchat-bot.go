package ipmibot

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"time"
	//"os"
)

type dbLogEntry struct {
	timestamp     int
	caller        string
	chatstring    string
	chatout       string
	systemcommand string
	systemout     string
	systemerror   int
}

type dbAliasEntry struct {
	name  string
	owner string
	host  string
}

var dbaddr string = "./ipmibot.sqlite3"
var DB *sql.DB = nil
var dbversion = 0

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
		log.Println("Database version match")
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
				systemerror int );`
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

/*
 timestamp
 caller
 chatstring
 chatout
 systemcommand
 systemout
 systemerror

*/

func logToDB(logEntry dbLogEntry) {
	statement, err := DB.Prepare(`
					insert into log(
					timestamp,
					caller,
					chatstring,
					chatout,
					systemcommand,
					systemout,
					systemerror) values (?, ?, ?, ?, ?, ?, ?);`)
	statement.Exec(
		int32(time.Now().Unix()),
		logEntry.caller,
		logEntry.chatstring,
		logEntry.chatout,
		logEntry.systemcommand,
		logEntry.systemout,
		logEntry.systemerror)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
}

func lastFromDB(params ...int) []dbLogEntry {
	nentries := 10
	if len(params) > 0 {
		nentries = params[0]
	}
	statement := fmt.Sprintf("select * from log order by timestamp desc limit %d", nentries)
	rows, err := DB.Query(statement)
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
		var systemerror int
		rows.Scan(&timestamp, &caller, &chatstring, &chatout, &systemcommand, &systemout, &systemerror)
		e := dbLogEntry{
			timestamp:     timestamp,
			caller:        caller,
			chatstring:    chatstring,
			chatout:       chatout,
			systemcommand: systemcommand,
			systemout:     systemout,
			systemerror:   systemerror,
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
	statement.Exec(alias.name, alias.owner)
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
	}
}

func showAlias(alias dbAliasEntry) []dbAliasEntry {
	var entries []dbAliasEntry
	statement := fmt.Sprintf("select * from hostaliases where name='%s' and owner='%s'", alias.name, alias.owner)
	rows, err := DB.Query(statement)
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

/*
		sqlStmt := `
	    create table foo (id integer not null primary key, name text);
	    delete from foo;
	    `
		_, err = db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Fatal(err)
		}
		stmt, err := tx.Prepare("insert into foo(id, name) values(?, ?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		for i := 0; i < 100; i++ {
			_, err = stmt.Exec(i, fmt.Sprintf("こんにちわ世界%03d", i))
			if err != nil {
				log.Fatal(err)
			}
		}
		tx.Commit()

		rows, err := db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}

		stmt, err = db.Prepare("select name from foo where id = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		var name string
		err = stmt.QueryRow("3").Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)

		_, err = db.Exec("delete from foo")
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
		if err != nil {
			log.Fatal(err)
		}

		rows, err = db.Query("select id, name from foo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			var name string
			err = rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(id, name)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
*/
