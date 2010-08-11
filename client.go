package ftp

import (
  "os"
)

type Connection struct {
  host string
  user string
  password string
}

func newConnection(host string, user string, password string) (*Connection) {
  return &Connection{host, user, password}
}

func Connect(host string, user string, password string) (*Connection, os.Error) {
  if host == "" {
    return nil, os.NewError("FTP Connection Error: Host can not be blank!")
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

