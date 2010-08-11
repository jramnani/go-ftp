package ftp

import (
    "fmt"
    /*"http"*/
)

type Connection struct {
  host string
  user string
  password string
}

/*type Error struct {*/
  /*Code int*/
  /*Message string*/
/*}*/

func (c *Connection) Connect() *Connection {
  fmt.Println("User ", c.user)
  fmt.Println("Pass", c.password)
  return c
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

