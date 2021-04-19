// main.go
package main

func main() {
  config := GetConfig()
  hosts := GetAllHostsFromSSHConfig(config)
  table := Format(hosts, config)
  // RenderTable(table, config)
  RenderFZF(table, hosts, config)
}
