package ftp

import (
  /*"fmt"*/
  "testing"
)

func TestConnect(t *testing.T) {
  _, err := Connect("ftp://localhost", "anonymous", "anonymous")
  if err != nil {
    t.Error(err)
  }
}

func TestConnectErrors(t *testing.T) {
  _, err := Connect("", "anonymous", "anonymous")
  if err == nil {
    t.Error("Connect() should have an error on blank 'host'!")
  }
  _, err = Connect("ftp://localhost", "", "anonymous")
  if err == nil {
    t.Error("Connect() should have an error on blank 'user'!")
  }
}
