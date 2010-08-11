package ftp

import (
  "fmt"
  "testing"
)

func TestConnect(t *testing.T) {
  conn := Connection{"ftp://localhost", "anonymous", "anonymous"}
  fmt.Println(conn)
}
