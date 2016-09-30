This is implementation of Hipchat integration bot to execute IPMI commands directly from
hipchat.

Valid commands are:
 - help - list of help topics, for more information type /ipmi help <topic>
 - reboot <ip or alias>
 - off <ip or alias>
 - on <ip or alias>
 - lanboot <ip or alias>
 - status <ip or alias>
 - alias [ - | add | del | show ] [<alias name>]
 - last [<number>]

For more details on Hipchat Integrations see [Hipchat documentation](https://www.hipchat.com/integrations)

# Configuration file
Configuration file is in YAML format, example configuration:

```yaml
address:
port:         8000
pidfile:      /var/run/hipchatbot/ipmi-hipchat-gobot.pid
logfile:      /var/log/hipchatbot/ipmi-hipchat-gobot.log
workdir:      .
ipmiusername: exampleUsername
ipmipassword: examplePassword
dbpath: /var/db/hipchatbot/ipmibot.sqlite3
```

# Startup
There is FreeBSD rc.d startup script:
etc/rc.d/ipmi-hipchat-gobot

# Dependencies
gopkg.in/yaml.v2
github.com/sevlyar/go-daemon
github.com/mattn/go-sqlite3

# Build
```
go get gopkg.in/yaml.v2
go get github.com/sevlyar/go-daemon
go get github.com/mattn/go-sqlite3
go test # recommended, but not mandatory
go build
```
