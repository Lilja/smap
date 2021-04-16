package main

import (
  "fmt"
  "strings"
  "log"
  "io/ioutil"
  "github.com/alexflint/go-arg"
)

var version = "0.1.0"

// Columns is a list of column.
type Columns []Column

// Config is the config of the program.
type Config struct {
  Verbose bool `default:"False"`
  Columns Columns `Help:"The columns to show. Edit me"`
  File string `help:"The ssh config file to read" default:"~/.ssh/config"`
}

// UnmarshalText is for go-arg for custom validating/parsing of the []Column from CLI
func (columns *Columns) UnmarshalText(b []byte) error {
  s := string(b)
  log.Println("Parsing columns")

  var err error
  if !strings.Contains(s, ",") {
    log.Println("Checking prop '", s, "'")
    column, _err := CheckColumnProperty(s)
    log.Println("Checking prop: error", err)
    if _err == nil {
      log.Println("Adding to column list", column)
      *columns = append(*columns, column)
    } else {
      log.Println("Adding error", _err)
      err = _err
    }
  } else {
    log.Println("[m] Checking prop for multiple")
    for _, prop := range strings.Split(s, ",") {
      log.Println("[m] Checking prop for multiple", prop)
      val, _err := CheckColumnProperty(prop)
      if _err == nil {
        log.Println("[m] Adding to column list", val)
        *columns = append(*columns, val)
      } else {
        log.Println("[m] Adding error", _err)
        err = _err
      }
    }
  }
  log.Println("Column Validator error", err)
  return err
}

// GetText returns the header to be used for table formatting of the specified columns
func (columns *Columns) GetText() []string {
  x := make([]string, len(*columns))
  for _, column := range *columns {
    x = append(x, columnValues[column])
  }
  return x
}

// GetConfig uses go-arg library and returns a config
func GetConfig() Config {
  var conf Config

  arg.MustParse(&conf)
  fmt.Println(conf.Verbose)
  fmt.Println(conf.Columns)

  if !conf.Verbose {
    log.SetFlags(0)
    log.SetOutput(ioutil.Discard)
  } else {
    log.Println("Verbose logging")
  }

  return conf
}

// Description is for the go-arg library.
func (Config) Description() string {
  return "Find hosts in ssh configs and uses fuzzy search to list these. Version " + version
}

// Version is for the go-arg library.
func (Config) Version() string {
  return "smap " + version
}
