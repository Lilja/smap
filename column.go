package main

import (
  "log"
  "errors"
  "strings"
)

// Column is a column that is printable in fzf
type Column int
const (
  // Host is the Host pattern in sshconfig
  Host = iota
  // Username is the Username directive in sshconfig
  Username
  // Hostname is the Hostname directive in sshconfig
  Hostname
  // Port is the Port directive in sshconfig
  Port
  // Jump is the ProxyJump directive in sshconfig
  Jump
  // NoJumps is the count of ProxyJumps, recursive.
  NoJumps
)

var columnValues = map[Column]string {
  Host: "Host",
  Username: "Username",
  Hostname: "Hostname",
  Port: "Port",
  Jump: "Jump",
  NoJumps: "Nojumps",
}
var columnKeys = map[string]Column {
  "Host": Host,
  "Username": Username,
  "Hostname": Hostname,
  "Port": Port,
  "Jump": Jump,
  "Nojumps": NoJumps,
}

func getColumnKeys() [] string {
  values := make([]string, 0, len(columnValues))

  for _, v := range columnValues {
    values = append(values, v)
  }
  return values
}

func safeKey(s string) string {
  return strings.Title(strings.ToLower(s))
}

// CheckColumnProperty checks if the user who submitted the column typed it correctly.
func CheckColumnProperty(s string) (Column, error) {
  safeK := safeKey(s)
  log.Println("Original", s, "Safe", safeK)
  v, safe := columnKeys[safeK]
  log.Println("Safe?", safe, "value", v)
  if safe {
    return v, nil
  }
  return v, errors.New("Unrecognized value '" + s + "'. Valid values are " + strings.Join(getColumnKeys(), ","))
}
