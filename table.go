package main

import (
  "log"
  "os"
  "os/exec"
  "strings"
  "github.com/olekukonko/tablewriter"
  "github.com/ktr0731/go-fuzzyfinder"
)

// renderTable pads and formats the data that is supposed to be printed with fzf
func renderTable(tableData [][]string, config Config) []string {
  tableString := &strings.Builder{}
  table := tablewriter.NewWriter(tableString)
  table.SetHeader(config.Columns.GetText())

  table.SetCenterSeparator("")
  table.SetColumnSeparator("")
  table.SetRowSeparator("")

  for _, v := range tableData {
      table.Append(v)
  }
  table.Render()
  return splitAndFilter(tableString.String())
}

func splitAndFilter(str string) []string {
  var output []string

  for idx, line := range strings.Split(str, "\n") {
    if idx == 1 {
      // Skip Table header
      continue
    } else if line == "" {
      // Skip Table formatted newline
      continue
    } else {
      output = append(output, line)
    }
  }

  return output
}

// RenderFZF takes data and renders with fzf
func RenderFZF(tableData [][]string, hosts []SSHHost, config Config) {
  formatedString := renderTable(tableData, config)
  idx, err := fuzzyfinder.Find(
      formatedString,
      func(i int) string {
          return formatedString[i]
      },
  )
  if err != nil {
    log.Fatal("Fuzzyfinder error ", err)
  } else {
    host := hosts[idx]
    log.Println("Selected index", idx, " Host: ", host, " name: ", host.name)
    cmd := exec.Command("ssh", host.name)
    cmd.Stdout = os.Stdout
    cmd.Stdin = os.Stdin
    cmd.Stderr = os.Stderr
    log.Println("Running cmd")
    e := cmd.Run()
    if e != nil {
      log.Fatal("Unable to run ssh command", e)
    }
  }
}
