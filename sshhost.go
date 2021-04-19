package main

// SSHHost is a representation of a Host pattern in a ssh config. With extra parameters that is custom made.
type SSHHost struct {
    name string
    user string
    hostname string
    port string
    proxyJump string
    noProxyJumps string
}
