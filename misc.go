package main

import (
	"fmt"
	"net"
)

func aliasToString(e dbAliasEntry) string {
	return fmt.Sprintf("'%s' is an alias for %s (owner %s)", e.name, e.host, e.owner)
}

func mkDbAliasEntry(name string, owner string, host string) *dbAliasEntry {
	var e dbAliasEntry
	if net.ParseIP(host) == nil {
		return nil
	}
	e.name, e.owner, e.host = name, owner, host
	return &e
}
