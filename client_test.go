package ftp

import (
  /*"fmt"*/
  "testing"
)

func TestDial(t *testing.T) {
  _, err := Dial("")
  if err == nil {
    t.Error("Dial() should have an error on blank 'host'!")
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
    t.Error("Login() should have an error on blank 'user'!")
  }
  err = conn.Login("anonymous", "")
  if err == nil {
    t.Error("Login() should have an error on blank 'password'!")
  }
}

