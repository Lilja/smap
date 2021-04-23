// main.go
package main

func main() {
  config := GetConfig()
  if config.CheckForUpdates {
    CheckForUpdates(config)
  } else {
    hosts := GetAllHostsFromSSHConfig(config)
    table := Format(hosts, config)
    RenderFZF(table, hosts, config)
  }
}
