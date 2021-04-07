  // main.go
  package main
  import (
      "github.com/kevinburke/ssh_config"
      "fmt"
      "os"
      "path/filepath"
  )
  func main() {
      f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
      cfg, _ := ssh_config.Decode(f)

      for _, host := range cfg.Hosts {
          fmt.Println("patterns:", host.Patterns)
          for _, node := range host.Nodes {
              // Manipulate the nodes as you see fit, or use a type switch to
              // distinguish between Empty, KV, and Include nodes.
              fmt.Println(node.String())
          }
      }
  }
