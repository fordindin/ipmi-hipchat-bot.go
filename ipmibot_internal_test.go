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

var _dbaddr string

var _ = json.Unmarshal(_jdata, &_jsonOut)

var _jsonDataTemplate = `{
    "event": "room_message",
    "item": {
        "message": {
            "date": "2016-09-06T21:53:37.420887+00:00",
            "from": {
                    "id": 14,
                    "links": {
                        "self": "https://hipchat.argotech.io/v2/user/14"
                    },
                    "mention_name": "TestUser",
                    "name": "Denis Barov",
                    "version": "BC5D9325"
            },
            "id": "f9f17833-9acd-48c8-9315-422e1b31e446",
            "mentions": [],
            "message": "%s",
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
}`

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
	config.Dbpath = trace()
	if DB != nil {
		_testsFailed += 1
		t.Error("DB should is not <nil> before initialization")
	}
	initDB()
	if DB == nil {
		_testsFailed += 1
		t.Error("DB shouldn't be <nil> after initialization")
	}
	os.Remove(config.Dbpath)
	//fmt.Println(reflect.TypeOf(DB))
}

func Test_closeDB(t *testing.T) {
	config.Dbpath = trace()
	initDB()
	closeDB()
	//fmt.Println(reflect.TypeOf(DB))
	os.Remove(config.Dbpath)
}

func Test_createDB(t *testing.T) {
	config.Dbpath = trace()
	initDB()
	//err = createDB()
	closeDB()
	//_ = err
	os.Remove(config.Dbpath)
}

func Test_verifyDB(t *testing.T) {
	config.Dbpath = trace()
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
	os.Remove(config.Dbpath)
}

func Test_verifyDB_nonexistent(t *testing.T) {
	config.Dbpath = trace()
	//initDB()
	//err := createDB()
	ver := verifyDB()
	if ver != -1 {
		_testsFailed += 1
		t.Error("Reporting version for nonexistent DB")
	}
	//_ = err
	_ = ver
}

func Test_logDB_lastFromDB_single(t *testing.T) {
	config.Dbpath = trace()
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
		//fmt.Println(entries[0])
		//fmt.Println(logEntry)
		_testsFailed += 1
		t.Error("Written log entry doesn't match")
	}
	closeDB()
	os.Remove(config.Dbpath)
	_testsPassed += 1
}

func Test_logDB_lastFromDB_many(t *testing.T) {
	config.Dbpath = trace()
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
	os.Remove(config.Dbpath)
	_testsPassed += 1
}

func Test_addAlias_showAlias(t *testing.T) {
	config.Dbpath = trace()
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
	os.Remove(config.Dbpath)
	_testsPassed += 1
}

func Test_addAlias_updateAlias(t *testing.T) {
	config.Dbpath = trace()
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
	os.Remove(config.Dbpath)
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
	var x = IpmiExec("10.20.30.13", "status")
	fmt.Printf("\n\tCommand:   %s\n", x[0].commandstring)
	fmt.Printf("\tOutput:    %s\n", x[0].output)
	_testsPassed += 1
}

func Test_IpmiExec_multi(t *testing.T) {
	trace()
	commands["testcommand"] = [][]string{
		[]string{"chassis", "power", "status"},
		[]string{"chassis", "power", "status"},
	}
	var x = IpmiExec("10.20.30.13", "testcommand")
	if len(x) != 2 {
		t.Error("There should be at least two returns in this command")
	} else {
		for _, o := range x {
			fmt.Printf("\n\tCommand:   %s\n", o.commandstring)
			fmt.Printf("\tOutput:    %s\n", o.output)
		}
	}
	_testsPassed += 1
}

func Test_parseInputJson(t *testing.T) {
	trace()
	s := "/ipmi some params here"
	_jsonData := fmt.Sprintf(_jsonDataTemplate, s)
	out := parseInputJson([]byte(_jsonData))
	_ = out
	if out.node != "/ipmi" ||
		out.name != "TestUser" ||
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
		_testsFailed += 1
		t.Error("Resulting alias entry doesn't match exemplary one")
	}
}

func Test_mkDbAliasEntry(t *testing.T) {
	trace()
	e := aliasEntry
	r := mkDbAliasEntry(e.name, e.owner, e.host)
	if r == nil {
		_testsFailed += 1
		t.Error("Resulting alias pointer shouldn't be nil")
	} else {
		ref := *r
		if aliasEntry != ref {
			_testsFailed += 1
			t.Error("Resulting alias entry doesn't match exemplary one")
		}
	}
}

func Test_mkDbAliasEntry_badAlias(t *testing.T) {
	trace()
	e := badAliasEntry
	if nil != mkDbAliasEntry(e.name, e.owner, e.host) {
		_testsFailed += 1
		t.Error("Resulting alias entry should be <nil>")
	}
}

func Test_processIpmi_alias(t *testing.T) {
	trace()
	initDB()
	exemplaryJsonStr := `{"color":"green","message":"@TestUser,\nAliases for you:\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		fmt.Println(out, exemplaryJsonStr)
		t.Error("Resulting output doesn't match exemplary one")
	}
	closeDB()
	os.Remove(config.Dbpath)
}

func Test_processIpmi_alias_add_bad_ip(t *testing.T) {
	trace()
	initDB()
	exemplaryJsonStr := `{"color":"red","message":"@TestUser,\n'notAnIPAddress' doesn't look like IP-address to me\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias add test notAnIPAddress"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
	closeDB()
	os.Remove(config.Dbpath)
}
func Test_processIpmi_alias_add(t *testing.T) {
	trace()
	initDB()
	exemplaryJsonStr := `{"color":"green","message":"@TestUser,\n'test' is an alias for 127.0.0.1 (owner TestUser)","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias add test 127.0.0.1"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
	closeDB()
	os.Remove(config.Dbpath)
}

func Test_processIpmi_alias_add_multi_than_list(t *testing.T) {
	trace()
	initDB()
	exemplaryJsonStr := `{"color":"green","message":"@TestUser,\nAliases for you:\n'test' is an alias for 127.0.0.1 (owner TestUser)\n'test2' is an alias for 10.0.0.1 (owner TestUser)\n'test3' is an alias for 172.16.0.1 (owner TestUser)\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias add test 127.0.0.1"))
	out := processIpmi(j)
	j = []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias add test2 10.0.0.1"))
	out = processIpmi(j)
	j = []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias add test3 172.16.0.1"))
	out = processIpmi(j)
	j = []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias"))
	out = processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
	//closeDB()
	//os.Remove(config.Dbpath)
	_dbaddr = config.Dbpath
}

func Test_processIpmi_alias_show_nonexistent(t *testing.T) {
	trace()
	exemplaryJsonStr := `{"color":"red","message":"@TestUser,\nThere is no alias 'test5' for TestUser\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias show test5"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_alias_show(t *testing.T) {
	trace()
	exemplaryJsonStr := `{"color":"green","message":"@TestUser,\n'test2' is an alias for 10.0.0.1 (owner TestUser)","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias show test2"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_alias_del(t *testing.T) {
	trace()
	exemplaryJsonStr := `{"color":"green","message":"@TestUser,\nAlias 'test2' removed (owner  TestUser)","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias del test2"))
	out := processIpmi(j)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_alias_del_nonexistent(t *testing.T) {
	trace()
	exemplaryJsonStr := `{"color":"red","message":"@TestUser,\nThere is no alias 'test7' for TestUser\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi alias del test7"))
	out := processIpmi(j)
	//fmt.Println(out)
	if out != exemplaryJsonStr {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_unwrapAlias_existent(t *testing.T) {
	trace()
	ret := unwrapAlias("test3", "TestUser")
	if ret != "172.16.0.1" {
		_testsFailed += 1
		t.Error("Resulting output doesn't match alias value")
	}
}

func Test_unwrapAlias_nonexistent(t *testing.T) {
	trace()
	ret := unwrapAlias("test11", "TestUser")
	if ret != "test11" {
		_testsFailed += 1
		t.Error("Resulting output doesn't match alias value")
	}
}

func Test_processIpmi_nonexistent(t *testing.T) {
	trace()
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi add"))
	out := processIpmi(j)
	exemplaryOutput := `{"color":"green","message":"@TestUser,\nValid commands are:\n\t\thelp - list of help topics, for more information type /ipmi help \u003ctopic\u003e\n\t\treboot \u003cip or alias\u003e\n\t\toff \u003cip or alias\u003e\n\t\ton \u003cip or alias\u003e\n\t\tlanboot \u003cip or alias\u003e\n\t\tstatus \u003cip or alias\u003e\n\t\talias [ - | add | del | show ] [\u003calias name\u003e]\n\t\tlast [\u003cnumber\u003e]","message_format":"text","notify":false}`
	if out != exemplaryOutput {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_last(t *testing.T) {
	trace()
	exemplaryOutput := `{"color":"green","message":"@TestUser,\nLast executed commands:\n 1     TestUser: /ipmi alias add test 127.0.0.1\n 2     TestUser: /ipmi alias add test2 10.0.0.1\n 3     TestUser: /ipmi alias add test3 172.16.0.1\n 4     TestUser: /ipmi alias \n 5     TestUser: /ipmi alias show test5\n 6     TestUser: /ipmi alias show test2\n 7     TestUser: /ipmi alias del test2\n 8     TestUser: /ipmi alias del test7\n 9     TestUser: /ipmi add \n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi last"))
	out := processIpmi(j)
	if out != exemplaryOutput {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_last_2(t *testing.T) {
	trace()
	exemplaryOutput := `{"color":"green","message":"@TestUser,\nLast executed commands:\n 1     TestUser: /ipmi alias add test 127.0.0.1\n 2     TestUser: /ipmi alias add test2 10.0.0.1\n","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi last 2"))
	out := processIpmi(j)
	if out != exemplaryOutput {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_processIpmi_last_notANumber(t *testing.T) {
	trace()
	exemplaryOutput := `{"color":"red","message":"@TestUser,\n'NaN' doesn't look like a number to me","message_format":"text","notify":false}`
	j := []byte(fmt.Sprintf(_jsonDataTemplate, "/ipmi last NaN"))
	out := processIpmi(j)
	if out != exemplaryOutput {
		_testsFailed += 1
		t.Error("Resulting output doesn't match exemplary one")
	}
}

func Test_readConfig(t *testing.T) {
	trace()
	config.Ipmiusername = "TestUsername"
	err := readConfig("etc/ipmi-hipchat-bot.cfg.example")
	if err != nil {
		_testsFailed += 1
		t.Error("Error reading default configuration file")
	}
	if config.Ipmiusername != "exampleUsername" {
		_testsFailed += 1
		t.Error("Configuration parameter override failed")
	}
}

func Test_readConfig_noConfigfile(t *testing.T) {
	trace()
	fmt.Println("Intended error messages:")
	err := readConfig("nonexistent.config.file")
	if err != nil {
		if config.Ipmiusername != "exampleUsername" {
			_testsFailed += 1
			t.Error("Default configuration settings failed")
		}
	} else {
		_testsFailed += 1
		t.Error("Should be nonexistent-file reading error")
	}
}

func Test_stat(t *testing.T) {
	closeDB()
	os.Remove(_dbaddr)
	stat := fmt.Sprintf(`
	Total Tests:% 6d
	Failures:% 9d
`,
		_testsTotal, _testsFailed)
	fmt.Println(stat)
}
