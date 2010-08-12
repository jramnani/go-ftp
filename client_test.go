package ftp

import (
  /*"fmt"*/
  "testing"
)

func TestConnect(t *testing.T) {
  _, err := Connect("localhost:21", "anonymous", "anonymous")
  if err != nil {
    t.Error(err)
  }
}

func TestConnectErrors(t *testing.T) {
  _, err := Connect("", "anonymous", "anonymous")
  if err == nil {
    t.Error("Connect() should have an error on blank 'host'!")
  }
  _, err = Connect("localhost:21", "", "anonymous")
  if err == nil {
    t.Error("Connect() should have an error on blank 'user'!")
  }
}

func TestLogin(t *testing.T) {
  conn, err := Connect("localhost:21", "anonymous", "anonymous")
  if err != nil {
    t.Fatal(err)
  }
  err = conn.Login()
  if err != nil {
    t.Error(err)
  }
}

