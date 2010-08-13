package ftp

import (
  "net"
  "os"
  "strings"
)

type Connection struct {
  host string
  user string
  password string
}

var CRLF = "\r\n"

func (c *Connection) Login() os.Error {
  conn, err := net.Dial("tcp", "", c.host)
  if err != nil {
    return err
  }
  defer conn.Close()
  _, err = conn.Write([]byte("USER " + c.user + CRLF))
  _, err = conn.Write([]byte("PASS " + c.password + CRLF))
  if err != nil {
    return err
  }
  return nil
}

func newConnection(host string, user string, password string) (*Connection) {
  return &Connection{host, user, password}
}

// Reused from src/pkg/http/client.go
// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

func Connect(host string, user string, password string) (*Connection, os.Error) {
  if host == "" {
    return nil, os.NewError("FTP Connection Error: Host can not be blank!")
  }
  if !hasPort(host) {
    return nil, os.NewError("FTP Connection Error: Host must have a port! e.g. host:21")
  }
  if user == "" {
    return nil, os.NewError("FTP Connection Error: User can not be blank!")
  }
  return newConnection(host, user, password), nil
}


func (c Connection) String() string {
  // Only print the password if the user is anonymous.
  if c.user == "anonymous" {
    return "FTP Connection " + "HOST: " + c.host + " USER: " + c.user + " PASS: " +c.password
  } else {
    return "FTP Connection " + "HOST: " + c.host + " USER: " + c.user + " PASS: [secret]"
  }
  return "The compiler made me do this."
}

