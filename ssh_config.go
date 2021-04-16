package main

import (
  "fmt"
  "log"
  "regexp"
  "os"
  "io/ioutil"
  "github.com/kevinburke/ssh_config"
  "strings"
  "path/filepath"
)

func getPath(filename string) (string) {
    m, _ := regexp.MatchString(`^/`, filename)
    if m {
      return filename
    }
    return filepath.Join(os.Getenv("HOME"), ".ssh", filename)
}

func _GetHosts(str string, allHosts []SSHHost) []SSHHost {
  cfg, err := ssh_config.Decode(strings.NewReader(str))
  if err != nil {
    log.Fatalln("Couldn't decode ssh config.", cfg)
  }

  includes := findAllIncludes(str)

  var allIncludes []SSHHost

  for _, host := range getHostNames(str) {
    newHost := createSSHHost(host, *cfg)
    allIncludes = append(allIncludes, *newHost)
  }

  for _, inc := range includes {
    content := readContentFromFile(getPath(inc))
    cfg, _ := ssh_config.Decode(strings.NewReader(content))

    var tempHosts []SSHHost
    for _, host := range getHostNames(content) {
      newHost := createSSHHost(host, *cfg)
      tempHosts = append(tempHosts, *newHost)
    }
    childHosts := _GetHosts(
      cfg.String(),
      append(allHosts, tempHosts...),
    )
    allIncludes = append(childHosts, allIncludes...)
  }

  return allIncludes
}

// expandTilde takes a path(str), and expands the home/tilde
func expandTilde(str string) string {
  if strings.HasPrefix(str, "~/") {
    return filepath.Join(os.Getenv("HOME"), str[2:])
  }
  return str
}


// GetAllHostsFromSSHConfig reads into the specified config file and creates an internal representation called SSHHost
func GetAllHostsFromSSHConfig(config Config) []SSHHost {
  var hosts []SSHHost
  filepath := expandTilde(config.File)
  log.Println("Reading file: ", config.File)
  content, err := ioutil.ReadFile(filepath)
  if err != nil {
    fmt.Println("Error")
    log.Fatal("ASDASDASD", err)
  }
  log.Println("Reading file: ", config.File)
  text := string(content)
  cfg, _ := ssh_config.Decode(strings.NewReader(text))
  return _GetHosts(cfg.String(), hosts)
}

func createSSHHost(host string, cfg ssh_config.Config) *SSHHost {
  a, _ := cfg.Get(host, "User")
  c, _ := cfg.Get(host, "Port")
  hostname, _ := cfg.Get(host, "Hostname")
  jump := getJump(cfg, host)

  sh := SSHHost{
    name: host,
    user: a,
    port: c,
    hostname: hostname,
    proxyJump: jump,
    noProxyJumps: "",
  }
  return &sh
}

func getHostNames(sshConfigFile string) []string {
  lines := strings.Split(sshConfigFile, "\n")
  pat := regexp.MustCompile(`^Host ([^*]*)$`)

  var hosts []string
  for _, line := range lines {
    hostsWithoutAsterix := pat.FindStringSubmatch(line)
    if len(hostsWithoutAsterix) > 0 {
      hosts = append(hosts, hostsWithoutAsterix[1])
    }
  }
  return hosts
}

func readContentFromFile(filepath string) string {
  content, _ := ioutil.ReadFile(filepath)
  return string(content)
}

func findAllIncludes(str string) []string {
  pat := regexp.MustCompile(`^Include (.*?)$`)
  var allIncludes []string
  for _, line := range strings.Split(str, "\n") {
    include := pat.FindStringSubmatch(line)
    if len(include) > 0 {
      allIncludes = append(allIncludes, include[1])
    }
  }
  return allIncludes
}

func getJump(cfg ssh_config.Config, host string) string {
  hasJump, err := cfg.Get(host, "ProxyJump")
  if err != nil {
    log.Fatalln("cfg.Get errored", err)
  }
  return hasJump
}
