package ftp

import (
  "fmt"
  "testing"
)

// TODO: Mock out network access.  Tests currently require an FTP server
// running on localhost.

func TestDial(t *testing.T) {
  _, err := Dial("")
  if err == nil {
    t.Error("Dial() should return an error on blank 'host'!")
  }
  _, err = Dial("localhost:21")
  if err != nil {
    t.Error(err)
  }
}

func TestLogin(t *testing.T) {
  conn, err := Dial("localhost:21")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login("anonymous", "anonymous@")
  if err != nil {
    t.Error(err)
  }
  err = conn.Login("", "anonymous@")
  if err == nil {
    t.Error("Login() should return an error on blank 'user'!")
  }
  err = conn.Login("anonymous", "")
  if err == nil {
    t.Error("Login() should return an error on blank 'password'!")
  }
}

func TestCmd(t *testing.T) {
  conn, err := Dial("localhost:ftp")
  if err != nil {
    t.Fatal(err)
  }
  code, response, err := conn.Cmd("USER", "anonymous")
  if err != nil {
    t.Error(err)
  }
  fmt.Println("Response Code:", code, "Response Text:", response)
}

