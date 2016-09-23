package main

import (
	"fmt"
	"os"
	//"reflect"
	"encoding/json"
	"errors"
	"path"
	"runtime"
	"testing"
)

var _testsPassed = 0
var _testsFailed = 0
var _testsTotal = 0

var _jsonOut interface{}
var _jdata = []byte(`{
		"color": "green",
		"message": "test",
		"notify": false,
		"message_format": "text"
}`)

var _ = json.Unmarshal(_jdata, &_jsonOut)

var _jsonData = []byte(`{
    "event": "room_message",
    "item": {
        "message": {
            "date": "2016-09-06T21:53:37.420887+00:00",
            "from": {
                    "id": 14,
                    "links": {
                        "self": "https://hipchat.argotech.io/v2/user/14"
                    },
                    "mention_name": "DenisBarov",
                    "name": "Denis Barov",
                    "version": "BC5D9325"
            },
            "id": "f9f17833-9acd-48c8-9315-422e1b31e446",
            "mentions": [],
            "message": "/chassis some params here",
            "type": "message"
        },
    "room": {
        "id": 8,
        "is_archived": false,
        "links": {
            "participants": "https://hipchat.argotech.io/v2/room/8/participant",
            "self": "https://hipchat.argotech.io/v2/room/8",
            "webhooks": "https://hipchat.argotech.io/v2/room/8/webhook"
        },
        "name": "Lab Stuff",
        "privacy": "public",
        "version": "XHEAE44M"
        }
    },
    "oauth_client_id": "40fdb3b6-bba9-4683-9946-67e3443beea6",
    "webhook_id": 11
}`)

var logEntry = dbLogEntry{
	caller:        "caller",
	chatstring:    "/some chat",
	chatout:       "return",
	systemcommand: "ls -al",
	systemout:     "..",
	systemerror:   errors.New("Some Error")}

var aliasEntry = dbAliasEntry{
	name:  "testing",
	owner: "dindin",
	host:  "127.0.0.1",
}

var badAliasEntry = dbAliasEntry{
	name:  "testing",
	owner: "dindin",
	host:  "not an IP address",
}

func trace() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[1])
	fmt.Printf("\t%s:%d %s\n", path.Base(file), line, path.Base(f.Name()))
	_testsTotal += 1
	return path.Base(fmt.Sprintf("%s.sqlite", path.Base(f.Name())))
}

func Test_initDB(t *testing.T) {
	dbaddr = trace()
	if DB != nil {
		_testsFailed += 1
		t.Error("DB should is not <nil> before initialization")
	}
	initDB()
	if DB == nil {
		_testsFailed += 1
		t.Error("DB shouldn't be <nil> after initialization")
	}
	os.Remove(dbaddr)
	//fmt.Println(reflect.TypeOf(DB))
	_testsPassed += 1
}

func Test_closeDB(t *testing.T) {
	dbaddr = trace()
	initDB()
	closeDB()
	//fmt.Println(reflect.TypeOf(DB))
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_createDB(t *testing.T) {
	dbaddr = trace()
	initDB()
	//err = createDB()
	closeDB()
	//_ = err
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_verifyDB(t *testing.T) {
	dbaddr = trace()
	dbversion = 65535
	initDB()
	//err := createDB()
	ver := verifyDB()
	if ver != dbversion {
		_testsFailed += 1
		t.Error("DB version mismatch")
	}
	closeDB()
	//_ = err
	_ = ver
	os.Remove(dbaddr)
	_testsPassed += 1
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
	_testsPassed += 1
}
*/

func Test_logDB_lastFromDB_single(t *testing.T) {
	dbaddr = trace()
	initDB()
	//err := createDB()
	//_ = err
	logToDB(logEntry)
	entries := lastFromDB()
	entries[0].timestamp = 0
	e1err := logEntry.systemerror
	e2err := entries[0].systemerror
	entries[0].systemerror = logEntry.systemerror
	if entries[0] != logEntry ||
		fmt.Sprintf("%s", e1err) != fmt.Sprintf("%s", e2err) {
		fmt.Println(entries[0])
		fmt.Println(logEntry)
		_testsFailed += 1
		t.Error("Written log entry doesn't match")
	}
	closeDB()
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_logDB_lastFromDB_many(t *testing.T) {
	dbaddr = trace()
	nentries := 20
	initDB()
	//err := createDB()
	//_ = err
	for i := 0; i < nentries; i++ {
		logToDB(logEntry)
	}
	entries := lastFromDB(nentries - 1)
	if len(entries) != nentries-1 {
		_testsFailed += 1
		t.Error("Requested number of entries doesn't match")
	}
	closeDB()
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_addAlias_showAlias(t *testing.T) {
	dbaddr = trace()
	initDB()
	//err := createDB()
	//_ = err
	addAlias(aliasEntry)
	aliases := showAlias(aliasEntry)
	if aliases[0] != aliasEntry {
		_testsFailed += 1
		t.Error("Created alias doesn't match template")
	}
	closeDB()
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_addAlias_updateAlias(t *testing.T) {
	dbaddr = trace()
	initDB()
	//err := createDB()
	//_ = err
	var e dbAliasEntry
	*(&e) = *(&aliasEntry)
	addAlias(e)
	e.host = "10.0.0.1"
	updateAlias(e)
	aliases := showAlias(e)
	if aliases[0] != e {
		_testsFailed += 1
		t.Error("Created alias doesn't match template")
	}
	closeDB()
	os.Remove(dbaddr)
	_testsPassed += 1
}

func Test_execCommand(t *testing.T) {
	trace()
	input := "some input"
	ret := execCommand("echo", input)
	if ret.output != input+"\n" {
		_testsFailed += 1
		t.Error("Input and Output are different")
	}
	if ret.err != nil {
		_testsFailed += 1
		t.Error("Command completed with error")
	}
	//t.Error("DB shouldn't be <nil> after initialization")
	_testsPassed += 1
}

func Test_IpmiExec(t *testing.T) {
	trace()
	if true {
		var x = IpmiExec("10.20.30.13", "status")
		fmt.Printf("\n\tCommand:   %s\n", x.commandstring)
		fmt.Printf("\tOutput:    %s\n", x.output)
	}
	_testsPassed += 1
}

func Test_parseInputJson(t *testing.T) {
	trace()
	out := parseInputJson(_jsonData)
	_ = out
	if out.node != "/chassis" ||
		out.name != "DenisBarov" ||
		out.command != "some" ||
		out.args[0] != "params" ||
		out.args[1] != "here" {
		_testsFailed += 1
		t.Error("Wrong field values in resulting json", out)
	}
	//fmt.Println(out)
	_testsPassed += 1
}

func Test_jsonFormatReply(t *testing.T) {
	trace()
	r := jsonFormatReply("green", "test")
	o, _ := json.Marshal(_jsonOut)
	if r != string(o) {
		_testsFailed += 1
		t.Error("Resulting json doesn't match with exemplary one")
	}
}

func Test_aliasToString(t *testing.T) {
	trace()
	e := aliasEntry
	var exemplaryAliasString = fmt.Sprintf("'%s' is an alias for %s (owner %s)", e.name, e.host, e.owner)
	if exemplaryAliasString != aliasToString(e) {
		t.Error("Resulting alias entry doesn't match exemplary one")
	}
}

func Test_mkDbAliasEntry(t *testing.T) {
	trace()
	e := aliasEntry
	r := mkDbAliasEntry(e.name, e.owner, e.host)
	if r == nil {
		t.Error("Resulting alias pointer shouldn't be nil")
	} else {
		ref := *r
		if aliasEntry != ref {
			t.Error("Resulting alias entry doesn't match exemplary one")
		}
	}
}

func Test_mkDbAliasEntry_badAlias(t *testing.T) {
	trace()
	e := badAliasEntry
	if nil != mkDbAliasEntry(e.name, e.owner, e.host) {
		t.Error("Resulting alias entry should be <nil>")
	}
}

func Test_stat(t *testing.T) {
	/*
			stat := fmt.Sprintf(`
				Tests total:	%d
				Tests passed:	%d
				Tests failed:	%d
		`,
				_testsTotal, _testsTotal, _testsFailed, (_testsTotal - _testsPassed - _testsFailed))
			fmt.Println(stat)
	*/
}
