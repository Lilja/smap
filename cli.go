package main

import (
  "os"
  "log"
  "io/ioutil"
  "github.com/alexflint/go-arg"
)

var version = "0.1.0"

// Columns is a list of column.
type Columns []Column

// Flags is what flags you can submit to the program
type Flags struct {
  Verbose bool `arg:"-v,--verbose" help:"verbosity level" default:"False"`
  Columns Columns `Help:"The columns to show. Edit me" default:"Host,Username,Port"`
  File string `help:"The ssh config file to read" default:"~/.ssh/config"`
  Update bool `arg:"--update" help:"Check for updates"`
}

// Config is the config of the program.
type Config struct {
  CurrentVersion string
  Verbose bool
  Columns Columns
  File string
  CheckForUpdates bool
}

// GetConfig uses go-arg library and returns a config
func GetConfig() Config {
  blockGoArg := false
  for _, x := range os.Args {
    if x == "-v" || x == "--verbose" {
      log.Println("Verbose logging, ")
      log.SetOutput(os.Stdout)
      blockGoArg = true
      break
    }
  }
  var flags Flags

  if !flags.Verbose && !blockGoArg {
    log.SetFlags(0)
    log.SetOutput(ioutil.Discard)
  }
  arg.MustParse(&flags)

  conf := Config {
    Verbose: flags.Verbose,
    CurrentVersion: version,
    Columns: flags.Columns,
    File: flags.File,
    CheckForUpdates: flags.Update,
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
