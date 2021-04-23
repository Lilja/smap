package main

import (
  "runtime"
  "os"
  "log"
  "io/ioutil"
  "github.com/alexflint/go-arg"
)

// These strings are changed with the github action goreleaser is changed with ldflags
var (
  version   = "dev"
  commit    = ""
  builtDate = ""
  builtBy   = ""
)

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
  version string
  builtDate string
  builtBy string
  commit string
}

// GetConfig uses go-arg library and returns a config
func GetConfig() Config {
  blockGoArg := false
  for _, x := range os.Args {
    if x == "-v" || x == "--verbose" {
      log.SetOutput(os.Stdout)
      log.Println("Verbose logging")
      log.Println("Device, os:", runtime.GOOS, "arch:", runtime.GOARCH)
      log.Println("Version", version)
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
    version: version,
    builtBy: builtBy,
    builtDate: builtDate,
    commit: commit,
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
