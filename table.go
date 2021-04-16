package main

import (
  "strings"
  "log"
  "github.com/olekukonko/tablewriter"
  "fmt"
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
func RenderFZF(tableData [][]string, config Config) {
  formatedString := renderTable(tableData, config)
  idx, err := fuzzyfinder.Find(
      formatedString,
      func(i int) string {
          return formatedString[i]
      })
  if err != nil {
      log.Fatal("FZF ", err)
  }
  fmt.Printf("selected: %v\n", idx)

}
