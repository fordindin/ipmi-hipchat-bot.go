package main

import (
	"fmt"
	"net"
)

// Converts dbAliasEntry to human-readable representation
func aliasToString(e dbAliasEntry) string {
	return fmt.Sprintf("'%s' is an alias for %s (owner %s)", e.name, e.host, e.owner)
}

// returns pointer to newly created dbAliasEntry based on input parameters
// additionaly checks validity of input IP-address. If an IP-address is invalid returns
// <nil> pointer.
func mkDbAliasEntry(name string, owner string, host string) *dbAliasEntry {
	var e dbAliasEntry
	if net.ParseIP(host) == nil {
		return nil
	}
	e.name, e.owner, e.host = name, owner, host
	return &e
}
