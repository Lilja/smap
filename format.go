package main

import (
	// "github.com/fatih/color"
  "log"
	"fmt"
	"sort"
)

func getAndFormatProperty(host SSHHost, hosts[]SSHHost, c Column) string {
  v := ""
  switch columnValues[c] {
    case "Host":
      // c := color.New(color.FgGreen).Add(color.Bold)
      // v = c.Sprint(host.name)
      v = host.name
    case "Hostname":
      v = host.hostname
    case "Username":
      v = host.user
    case "Port":
      v = host.port
    case "Jump":
      v = host.proxyJump
    case "Nojumps":
      log.Println("Calculating no jumps for", host.name)
      v = calculateNoJumps(host, hosts)
    default:
      panic(columnValues[c] + " case not implemented in switch case in getAndFormatProperty")
 }
 return v
}


func formatHost(host SSHHost, hosts []SSHHost, config Config) []string {
  // c := color.New(color.FgGreen).Add(color.Bold)
  // hostname := color.New(color.FgBlue).Add(color.Bold)

  row := make([]string, len(config.Columns))
  for _, c := range config.Columns {
      row = append(row, getAndFormatProperty(host, hosts, c))
  }

  return row
}

func recursiveJumpCount(hostName string, hosts []SSHHost, count int) int {
  if count >= 10 {
    log.Println("Panicing. Hostname", hostName)
    panic("Recursive deeper than 10 levels")
  }
  for _, x := range hosts {
    if x.name == hostName {
      if len(x.proxyJump) > 0 {
        log.Println("Finding proxyjump for", x.name, "depth=", count, x.proxyJump)
        return recursiveJumpCount(x.proxyJump, hosts, count + 1)
      }
      return count + 1
    }
  }
  return count
}

func calculateNoJumps(host SSHHost, hosts []SSHHost) string {
  c := recursiveJumpCount(host.proxyJump, hosts, 0)
  switch c {
    case 0:
      return "-"
    case 1:
      return "1 jump"
    default:
      return fmt.Sprintf("%d %s", c, "jumps")
  }
}

// Format formats according to the rules of Config
func Format(hosts []SSHHost, config Config) [][]string {
  sort.SliceStable(hosts, func(i, j int) bool {
    return hosts[i].hostname < hosts[j].hostname
  })
  var tableData [][]string
  for _, host := range hosts {
    tableData = append(tableData, formatHost(host, hosts, config))
  }

  return tableData
}
