package ipmibot

import (
	"fmt"
	"os"
	//"reflect"
	"path"
	"runtime"
	"testing"
)

var _dbaddr string = "file.sqlite3"

var logEntry = dbLogEntry{
	caller:        "caller",
	chatstring:    "/some chat",
	chatout:       "return",
	systemcommand: "ls -al",
	systemout:     "..",
	systemerror:   0}

var aliasEntry = dbAliasEntry{
	name:  "testing",
	owner: "dindin",
	host:  "127.0.0.1",
}

func trace() {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[1])
	fmt.Printf("\t%s:%d %s\n", path.Base(file), line, path.Base(f.Name()))
}

func Test_initDB(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	if DB != nil {
		t.Error("DB should is not <nil> before initialization")
	}
	initDB()
	if DB == nil {
		t.Error("DB shouldn't be <nil> after initialization")
	}
	os.Remove(_dbaddr)
	//fmt.Println(reflect.TypeOf(DB))
}

func Test_closeDB(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	closeDB()
	//fmt.Println(reflect.TypeOf(DB))
	os.Remove(_dbaddr)
}

func Test_createDB(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	//err = createDB()
	closeDB()
	//_ = err
	os.Remove(_dbaddr)
}

func Test_verifyDB(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	dbversion = 65535
	initDB()
	//err := createDB()
	ver := verifyDB()
	if ver != dbversion {
		t.Error("DB version mismatch")
	}
	closeDB()
	//_ = err
	_ = ver
	os.Remove(_dbaddr)
}

/*
func Test_verifyDB_nonexistent(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	//err := createDB()
	ver := verifyDB()
	if ver != -1 {
		t.Error("Reporting version for nonexistent DB")
	}
	closeDB()
	//_ = err
	_ = ver
	os.Remove(_dbaddr)
}
*/

func Test_logDB_lastFromDB_single(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	//err := createDB()
	//_ = err
	logToDB(logEntry)
	entries := lastFromDB()
	entries[0].timestamp = 0
	if entries[0] != logEntry {
		fmt.Println(entries[0])
		fmt.Println(logEntry)
		t.Error("Written log entry doesn't match")
	}
	closeDB()
	os.Remove(_dbaddr)
}

func Test_logDB_lastFromDB_many(t *testing.T) {
	trace()
	nentries := 20
	dbaddr = _dbaddr
	initDB()
	//err := createDB()
	//_ = err
	for i := 0; i < nentries; i++ {
		logToDB(logEntry)
	}
	entries := lastFromDB(nentries - 1)
	if len(entries) != nentries-1 {
		t.Error("Requested number of entries doesn't match")
	}
	closeDB()
	os.Remove(_dbaddr)
}

func Test_addAlias_showAlias(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	//err := createDB()
	//_ = err
	addAlias(aliasEntry)
	aliases := showAlias(aliasEntry)
	if aliases[0] != aliasEntry {
		t.Error("Created alias doesn't match template")
	}
	closeDB()
	os.Remove(_dbaddr)
}

func Test_addAlias_updateAlias(t *testing.T) {
	trace()
	dbaddr = _dbaddr
	initDB()
	//err := createDB()
	//_ = err
	addAlias(aliasEntry)
	aliasEntry.host = "10.0.0.1"
	updateAlias(aliasEntry)
	aliases := showAlias(aliasEntry)
	if aliases[0] != aliasEntry {
		t.Error("Created alias doesn't match template")
	}
	closeDB()
	os.Remove(_dbaddr)
}

func Test_execCommand(t *testing.T) {
	trace()
	input := "some input"
	ret := execCommand("echo", input)
	if ret.output != input+"\n" {
		t.Error("Input and Output are different")
	}
	if ret.err != nil {
		t.Error("Command completed with error")
	}
	//t.Error("DB shouldn't be <nil> after initialization")
}
