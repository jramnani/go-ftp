package ftp

import (
  "io"
  "net"
  "os"
  "strings"
)

// Knows the control connection where commands are sent to the server.
type Connection struct {
  control io.ReadWriteCloser
}

var CRLF = "\r\n"

func (c *Connection) cmd() {
}

// Log into a FTP server using username and password.
func (c *Connection) Login(user string, password string) os.Error {
  if user == "" {
    return os.NewError("FTP Connection Error: User can not be blank!")
  }
  if password == "" {
    return os.NewError("FTP Connection Error: Password can not be blank!")
  }
  _, err := c.control.Write([]byte("USER " + user + CRLF))
  _, err = c.control.Write([]byte("PASS " + password + CRLF))
  if err != nil {
    return err
  }
  return nil
}

// Dials up a remote FTP server.
// host should be in the form of address:port e.g. myserver:21 or myserver:ftp
func Dial(host string) (*Connection, os.Error) {
  if host == "" {
    return nil, os.NewError("FTP Connection Error: Host can not be blank!")
  }
  if !hasPort(host) {
    return nil, os.NewError("FTP Connection Error: Host must have a port! e.g. host:21")
  }
  conn, err := net.Dial("tcp", "", host)
  if err != nil {
    return nil, err
  }
  return &Connection{conn}, nil
}


// Reused from src/pkg/http/client.go
// Given a string of the form "host", "host:port", or "[ipv6::address]:port",
// return true if the string includes a port.
func hasPort(s string) bool { return strings.LastIndex(s, ":") > strings.LastIndex(s, "]") }

